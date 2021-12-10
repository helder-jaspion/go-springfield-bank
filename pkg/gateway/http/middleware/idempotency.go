package middleware

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/hlog"

	"github.com/helder-jaspion/go-springfield-bank/pkg/appcontext"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

const (
	cacheDur               = 24 * time.Hour
	headerIdempotencyKey   = "X-Idempotency-Key"
	headerIdempotencyCache = "X-Idempotency-Cache"
	cacheHit               = "HIT"
)

type response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

func generateHashKey(r *http.Request) string {
	idempotencyKey := r.Header.Get(headerIdempotencyKey)
	if idempotencyKey == "" {
		return ""
	}

	sub, _ := appcontext.GetAuthSubject(r.Context())
	hashKeyBytes := sha1.Sum([]byte(sub + "." + idempotencyKey + "." + r.Method + "." + r.RequestURI))
	return hex.EncodeToString(hashKeyBytes[:])
}

func saveResponse(ctx context.Context, idpRepo repository.IdempotencyRepository, rec *httptest.ResponseRecorder, hashKey string) (*response, error) {
	resp := &response{
		StatusCode: rec.Code,
		Headers:    rec.Header(),
		Body:       rec.Body.Bytes(),
	}

	content, err := json.Marshal(resp)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal response to json")
	}

	err = idpRepo.Set(ctx, hashKey, content, cacheDur)
	if err != nil {
		return nil, errors.Wrap(err, "could not cache response")
	}

	return resp, nil
}

// Idempotency returns the same result for requests with the same uri, user and X-Idempotency-Key header.
//
// Fallbacks to original request processing in case of errors.
func Idempotency(idpRepo repository.IdempotencyRepository, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := hlog.FromRequest(r)

		hashKey := generateHashKey(r)
		if hashKey == "" {
			next(w, r)
			return
		}

		respBytes, err := idpRepo.Get(r.Context(), hashKey)
		var resp *response
		if err == nil {
			err = json.Unmarshal(respBytes, &resp)
			if err != nil {
				logger.Error().Err(err).Interface("resp", resp).Msg("Could not marshal response to json.")
				next(w, r)
				return
			}
			resp.Headers.Add(headerIdempotencyCache, cacheHit)
		} else {
			rec := httptest.NewRecorder()
			next(rec, r)

			resp, err = saveResponse(r.Context(), idpRepo, rec, hashKey)
			if err != nil {
				logger.Error().Err(err).Interface("resp", resp).Msg("Could not cache response.")
				next(w, r)
				return
			}
		}

		for k, v := range resp.Headers {
			w.Header()[k] = v
		}

		w.WriteHeader(resp.StatusCode)
		_, err = w.Write(resp.Body)
		if err != nil {
			logger.Error().Err(err).Interface("resp", resp).Msg("Could not write response.")
			next(w, r)
			return
		}
	}
}

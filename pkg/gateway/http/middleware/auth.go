package middleware

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/hlog"

	"github.com/helder-jaspion/go-springfield-bank/pkg/appcontext"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/io"
)

// BearerAuth get Bearer Authorization header, parses and validate.
func BearerAuth(authUC usecase.AuthUseCase, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := hlog.FromRequest(r)

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			logger.Warn().Str("Authorization", r.Header.Get("Authorization")).Msg("malformed token")
			io.WriteErrorMsg(w, logger, http.StatusUnauthorized, "malformed Token")
			return
		}

		tokenClaims, err := authUC.Authorize(logger.WithContext(r.Context()), authHeader[1])
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err == usecase.ErrAuthInvalidAccessToken {
				statusCode = http.StatusUnauthorized
			}

			io.WriteErrorMsg(w, logger, statusCode, err.Error())
			return
		}

		next(w, r.WithContext(appcontext.WithAuthSubject(r.Context(), tokenClaims.Subject)))
	}
}

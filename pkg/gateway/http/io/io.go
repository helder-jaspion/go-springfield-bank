package io

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

const (
	contentType     = "Content-Type"
	jsonContentType = "application/json"
)

// ErrorOutput represents the output data in case of error.
type ErrorOutput struct {
	Code    int    `json:"code"`
	Message string `json:"message" example:"something wrong happened"`
}

// ReadInput reads the JSON-encoded value from request and stores it in the value pointed to by value.
func ReadInput(r *http.Request, logger *zerolog.Logger, value interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Error().Stack().Err(err).Msg("error closing request body")
		}
	}()

	return nil
}

// WriteSuccess writes a success result to the http.ResponseWriter
func WriteSuccess(w http.ResponseWriter, logger *zerolog.Logger, statusCode int, result interface{}) {
	w.Header().Set(contentType, jsonContentType)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		logger.Error().Stack().Err(err).Interface("result", result).Msg("error encoding response")
	}
}

// WriteErrorMsg writes an error message to the http.ResponseWriter
func WriteErrorMsg(w http.ResponseWriter, logger *zerolog.Logger, statusCode int, message string) {
	w.Header().Set(contentType, jsonContentType)
	w.WriteHeader(statusCode)

	errReturn := ErrorOutput{
		Code:    statusCode,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(errReturn); err != nil {
		logger.Error().Stack().Err(err).Msg("error encoding response")
	}
}

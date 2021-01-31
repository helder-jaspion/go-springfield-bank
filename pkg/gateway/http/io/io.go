package io

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"net/http"
)

// ReadInput reads the JSON-encoded value from request and stores it in the value pointed to by value.
func ReadInput(r *http.Request, logger *zerolog.Logger, value interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Error().Err(err).Msg("error closing request body")
		}
	}()

	return nil
}

// WriteError writes an error message to the http.ResponseWriter
func WriteError(w http.ResponseWriter, logger *zerolog.Logger, statusCode int, message string) {
	w.WriteHeader(statusCode)

	// TODO formato do retorno {code, message}
	errReturn := make(map[string]interface{})
	errReturn["code"] = statusCode
	errReturn["message"] = message

	if err := json.NewEncoder(w).Encode(errReturn); err != nil {
		logger.Error().Err(err).Msg("error encoding response")
	}
}

// WriteSuccess writes a success result to the http.ResponseWriter
func WriteSuccess(w http.ResponseWriter, logger *zerolog.Logger, statusCode int, result interface{}) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		logger.Error().Err(err).Interface("result", result).Msg("error encoding response")
	}
}

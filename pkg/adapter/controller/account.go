package controller

import (
	"encoding/json"
	"fmt"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/rs/zerolog/hlog"
	"net/http"
)

// AccountController is the interface that wraps http handle methods related to accounts
type AccountController interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type accountController struct {
	accountUC usecase.AccountUseCase
}

//NewAccountController instantiates a new account controller
func NewAccountController(accountUC usecase.AccountUseCase) AccountController {
	return &accountController{
		accountUC: accountUC,
	}
}

func (a accountController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var input usecase.AccountCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		hlog.FromRequest(r).Error().Err(err).Msg("error decoding account create input")
		writeError(w, http.StatusBadRequest, "error reading input")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			hlog.FromRequest(r).Error().Err(err).Msg("error closing request body")
		}
	}()

	ctx := hlog.FromRequest(r).WithContext(r.Context())
	result, err := a.accountUC.Create(ctx, input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		hlog.FromRequest(r).Error().Err(err).Msg("error encoding account create response")
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	// TODO formato do retorno {code, message}
	errReturn := make(map[string]interface{})
	errReturn["code"] = statusCode
	errReturn["message"] = message
	if err := json.NewEncoder(w).Encode(errReturn); err != nil {
		fmt.Printf("error encoding response: %v", err)
	}
}

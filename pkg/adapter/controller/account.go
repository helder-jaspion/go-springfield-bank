package controller

import (
	"encoding/json"
	"fmt"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
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
		writeError(w, http.StatusBadRequest, "error reading input")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			fmt.Printf("error closing request body: %v", err)
		}
	}()

	result, err := a.accountUC.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error creating account")
		return
	}
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		fmt.Printf("error encoding account create response: %v", err)
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

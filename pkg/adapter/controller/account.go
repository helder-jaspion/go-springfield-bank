package controller

import (
	"encoding/json"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"net/http"
)

// AccountController is the interface that wraps http handle methods related to accounts
type AccountController interface {
	Create(resp http.ResponseWriter, req *http.Request)
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

// Create handles account creation.
func (a accountController) Create(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-type", "application/json")
	var input usecase.AccountCreateInput
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		//json.NewEncoder(resp).Encode(errors.ServiceError{Message: "Error unmarshalling the request"})
		return
	}

	result, err := a.accountUC.Create(req.Context(), input)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		//json.NewEncoder(resp).Encode(errors.ServiceError{Message: "Error saving the input"})
		return
	}
	resp.WriteHeader(http.StatusCreated)
	json.NewEncoder(resp).Encode(result)
}

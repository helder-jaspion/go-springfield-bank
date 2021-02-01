package controller

import (
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/io"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"net/http"
)

// AccountController is the interface that wraps http handle methods related to accounts
type AccountController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
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
	logger := hlog.FromRequest(r)

	var input usecase.AccountCreateInput
	if err := io.ReadInput(r, logger, &input); err != nil {
		logger.Error().Err(err).Msg("error decoding account create input")
		io.WriteErrorMsg(w, logger, http.StatusBadRequest, "error reading input")
		return
	}

	result, err := a.accountUC.Create(logger.WithContext(r.Context()), input)
	if err != nil {
		a.writeError(w, logger, http.StatusInternalServerError, err)
		return
	}

	io.WriteSuccess(w, logger, http.StatusCreated, result)
}

func (a accountController) Fetch(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	result, err := a.accountUC.Fetch(logger.WithContext(r.Context()))
	if err != nil {
		a.writeError(w, logger, http.StatusInternalServerError, err)
		return
	}

	io.WriteSuccess(w, logger, http.StatusOK, result)
}

func (a *accountController) GetBalance(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	params := httprouter.ParamsFromContext(r.Context())

	result, err := a.accountUC.GetBalance(r.Context(), model.AccountID(params.ByName("id")))
	if err != nil {
		a.writeError(w, logger, http.StatusInternalServerError, err)
		return
	}

	io.WriteSuccess(w, logger, http.StatusOK, result)
}

func (a *accountController) writeError(w http.ResponseWriter, logger *zerolog.Logger, statusCode int, err error) {
	switch err {
	case repository.ErrAccountNotFound:
		statusCode = http.StatusNotFound
	case usecase.ErrAccountCPFAlreadyExists:
		statusCode = http.StatusConflict
	case usecase.ErrAccountNameWrongLength,
		usecase.ErrAccountSecretWrongLength,
		usecase.ErrAccountBalanceNegative,
		usecase.ErrAccountCPFInvalid:
		statusCode = http.StatusBadRequest
	}

	io.WriteErrorMsg(w, logger, statusCode, err.Error())
}

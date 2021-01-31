package controller

import (
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/io"
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
	logger := hlog.FromRequest(r)
	w.Header().Set("Content-type", "application/json")

	var input usecase.AccountCreateInput
	if err := io.ReadInput(r, logger, &input); err != nil {
		logger.Error().Err(err).Msg("error decoding account create input")
		io.WriteError(w, logger, http.StatusBadRequest, "error reading input")
		return
	}

	ctx := logger.WithContext(r.Context())
	result, err := a.accountUC.Create(ctx, input)
	if err != nil {
		io.WriteError(w, logger, http.StatusInternalServerError, err.Error())
		return
	}

	io.WriteSuccess(w, logger, http.StatusCreated, result)
}

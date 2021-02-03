package controller

import (
	"github.com/helder-jaspion/go-springfield-bank/pkg/appcontext"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/io"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"net/http"
)

// TransferController is the interface that wraps http handle methods related to transfers.
type TransferController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
}

type transferController struct {
	trfUC  usecase.TransferUseCase
	authUC usecase.AuthUseCase
}

//NewTransferController instantiates a new transfer controller.
func NewTransferController(trfUC usecase.TransferUseCase, authUC usecase.AuthUseCase) TransferController {
	return &transferController{
		trfUC:  trfUC,
		authUC: authUC,
	}
}

func (trfCtrl transferController) Create(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	accountID, ok := appcontext.GetAuthSubject(r.Context())
	if !ok {
		trfCtrl.writeError(w, logger, http.StatusUnauthorized, usecase.ErrAuthInvalidAccessToken)
		return
	}

	var input usecase.TransferCreateInput
	if err := io.ReadInput(r, logger, &input); err != nil {
		logger.Error().Err(err).Msg("error decoding transfer create input")
		io.WriteErrorMsg(w, logger, http.StatusBadRequest, "error reading input")
		return
	}
	input.AccountOriginID = accountID

	result, err := trfCtrl.trfUC.Create(logger.WithContext(r.Context()), input)
	if err != nil {
		trfCtrl.writeError(w, logger, http.StatusInternalServerError, err)
		return
	}

	io.WriteSuccess(w, logger, http.StatusCreated, result)
}

func (trfCtrl transferController) Fetch(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	accountID, ok := appcontext.GetAuthSubject(r.Context())
	if !ok {
		trfCtrl.writeError(w, logger, http.StatusUnauthorized, usecase.ErrAuthInvalidAccessToken)
		return
	}

	result, err := trfCtrl.trfUC.Fetch(logger.WithContext(r.Context()), model.AccountID(accountID))
	if err != nil {
		trfCtrl.writeError(w, logger, http.StatusInternalServerError, err)
		return
	}

	io.WriteSuccess(w, logger, http.StatusOK, result)
}

func (trfCtrl transferController) writeError(w http.ResponseWriter, logger *zerolog.Logger, statusCode int, err error) {
	switch err {
	case repository.ErrAccountNotFound,
		usecase.ErrAccountCurrentBalanceInsufficient:
		statusCode = http.StatusUnprocessableEntity
	case usecase.ErrTransferOriginAccountRequired,
		usecase.ErrTransferDestinationAccountRequired,
		usecase.ErrTransferAmountNotPositive,
		usecase.ErrTransferSameAccount:
		statusCode = http.StatusBadRequest
	case usecase.ErrAuthInvalidAccessToken:
		statusCode = http.StatusUnauthorized
	}

	io.WriteErrorMsg(w, logger, statusCode, err.Error())
}

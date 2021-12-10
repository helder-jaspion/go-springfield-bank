package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/io"
)

// AccountController is the interface that wraps http handle methods related to the accounts.
type AccountController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
}

type accountController struct {
	accUC usecase.AccountUseCase
}

//NewAccountController instantiates a new account controller.
func NewAccountController(accUC usecase.AccountUseCase) AccountController {
	return &accountController{
		accUC: accUC,
	}
}

// @Summary Create account
// @Description Creates a new account
// @tags Accounts
// @Accept  json
// @Produce  json
// @Param account body usecase.AccountCreateInput true "Account"
// @Param X-Idempotency-Key header string false "Idempotency key"
// @Success 201 {object} usecase.AccountCreateOutput
// @failure 400 {object} io.ErrorOutput
// @failure 409 {object} io.ErrorOutput
// @failure 500 {object} io.ErrorOutput
// @Router /accounts [post]
func (accCtrl accountController) Create(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	var input usecase.AccountCreateInput
	if err := io.ReadInput(r, logger, &input); err != nil {
		logger.Error().Stack().Err(err).Msg("error decoding account create input")
		io.WriteErrorMsg(w, logger, http.StatusBadRequest, "error reading input")
		return
	}

	result, err := accCtrl.accUC.Create(logger.WithContext(r.Context()), input)
	if err != nil {
		accCtrl.writeError(w, logger, http.StatusInternalServerError, err)
		return
	}

	io.WriteSuccess(w, logger, http.StatusCreated, result)
}

// @Summary Fetch accounts
// @Description Fetch all the accounts
// @tags Accounts
// @Produce json
// @Success 200 {object} []usecase.AccountFetchOutput
// @failure 500 {object} io.ErrorOutput
// @Router /accounts [get]
func (accCtrl accountController) Fetch(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	result, err := accCtrl.accUC.Fetch(logger.WithContext(r.Context()))
	if err != nil {
		accCtrl.writeError(w, logger, http.StatusInternalServerError, err)
		return
	}

	io.WriteSuccess(w, logger, http.StatusOK, result)
}

// @Summary Get account balance
// @Description Get the balance of an account
// @tags Accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} usecase.AccountBalanceOutput
// @failure 404 {object} io.ErrorOutput
// @failure 500 {object} io.ErrorOutput
// @Router /account/{id}/balance [get]
func (accCtrl accountController) GetBalance(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	params := httprouter.ParamsFromContext(r.Context())

	result, err := accCtrl.accUC.GetBalance(r.Context(), model.AccountID(params.ByName("id")))
	if err != nil {
		accCtrl.writeError(w, logger, http.StatusInternalServerError, err)
		return
	}

	io.WriteSuccess(w, logger, http.StatusOK, result)
}

func (accCtrl accountController) writeError(w http.ResponseWriter, logger *zerolog.Logger, statusCode int, err error) {
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

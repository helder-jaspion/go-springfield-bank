package controller

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/http/io"
)

// AuthController is the interface that wraps http handle methods related to authentication.
type AuthController interface {
	Login(w http.ResponseWriter, r *http.Request)
	//Logout(w http.ResponseWriter, r *http.Request)
	//RefreshToken(w http.ResponseWriter, r *http.Request)
}

type authController struct {
	authUC usecase.AuthUseCase
}

//NewAuthController instantiates a new auth controller.
func NewAuthController(authUC usecase.AuthUseCase) AuthController {
	return &authController{
		authUC: authUC,
	}
}

// @Summary Login
// @Description Authenticates the user/account
// @tags Authentication
// @Accept json
// @Produce json
// @Param credentials body usecase.AuthLoginInput true "Credentials"
// @Success 200 {object} usecase.AuthTokenOutput
// @failure 400 {object} io.ErrorOutput
// @failure 401 {object} io.ErrorOutput
// @Router /login [post]
func (authCtrl authController) Login(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	var input usecase.AuthLoginInput
	if err := io.ReadInput(r, logger, &input); err != nil {
		logger.Error().Stack().Err(err).Msg("error decoding login input")
		io.WriteErrorMsg(w, logger, http.StatusBadRequest, "error reading input")
		return
	}

	result, err := authCtrl.authUC.Login(logger.WithContext(r.Context()), input)
	if err != nil {
		authCtrl.writeError(w, logger, http.StatusInternalServerError, err)
		return
	}

	io.WriteSuccess(w, logger, http.StatusOK, result)
}

func (authCtrl authController) writeError(w http.ResponseWriter, logger *zerolog.Logger, statusCode int, err error) {
	switch err {
	case usecase.ErrAuthInvalidCredentials:
		statusCode = http.StatusUnauthorized
	}

	io.WriteErrorMsg(w, logger, statusCode, err.Error())
}

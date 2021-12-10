package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

var (
	// ErrAuthInvalidCredentials happens if the credentials are not recognized as valid.
	ErrAuthInvalidCredentials = errors.New("invalid credentials")
	// ErrAuthLogin happens when an error occurred while processing login.
	ErrAuthLogin = errors.New("could not login")
)

// AuthLoginInput represents the expected input data when logging in.
type AuthLoginInput struct {
	CPF    string `json:"cpf" example:"999.999.999-99"`
	Secret string `json:"secret" example:"S3cr3t"`
}

// AuthTokenOutput represents the output data of the login method.
type AuthTokenOutput struct {
	AccessToken string `json:"access_token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTI1Mjk3NTksImlhdCI6MTYxMjUyODg1OSwic3ViIjoiNmMzYjhhNTUtNmI4MC00MTM3LTlkZmYtNTAzY2FmNTc2NTE0In0.uJc24K1GLMwbswa_4DqCpNzqhdLYAklVyGaZEcmlKn8"`
}

func newAuthTokenOutput(accessToken string) *AuthTokenOutput {
	return &AuthTokenOutput{
		AccessToken: accessToken,
	}
}

// Login checks if the user credentials are valid and, if valid, returns a jwt access token.
func (authUC authUseCase) Login(ctx context.Context, loginInput AuthLoginInput) (*AuthTokenOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cpf := model.NewCPF(loginInput.CPF)

	account, err := authUC.accRepo.GetByCPF(ctx, cpf)
	if err != nil {
		if err == repository.ErrAccountNotFound {
			return nil, ErrAuthInvalidCredentials
		}
		log.Ctx(ctx).Error().Stack().Err(err).Str("cpf", loginInput.CPF).Msg("error during login")
		return nil, ErrAuthLogin
	}

	if account == nil || account.ID == "" {
		return nil, ErrAuthInvalidCredentials
	}

	err = account.CompareSecrets(loginInput.Secret)
	if err != nil {
		return nil, ErrAuthInvalidCredentials
	}

	authTokenOutput, err := authUC.createAccountToken(account.ID)
	if err != nil {
		log.Ctx(ctx).Error().Stack().Err(err).Str("cpf", loginInput.CPF).Str("accountID", string(account.ID)).Msg("error creating new authTokenOutput")
		return nil, ErrAuthLogin
	}

	return authTokenOutput, nil
}

func (authUC authUseCase) createAccountToken(accountID model.AccountID) (*AuthTokenOutput, error) {
	now := time.Now()
	accessTokenClaims := jwt.StandardClaims{
		Subject:   string(accountID),
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(authUC.accessTokenDur).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(authUC.secretKey))
	if err != nil {
		return nil, err
	}

	return newAuthTokenOutput(accessTokenString), nil
}

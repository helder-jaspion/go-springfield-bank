package mock

import (
	"context"

	"github.com/golang-jwt/jwt/v4"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
)

// AuthUseCase mocks an usecase.AuthUseCase.
type AuthUseCase struct {
	OnLogin     func(ctx context.Context, loginInput usecase.AuthLoginInput) (*usecase.AuthTokenOutput, error)
	OnAuthorize func(ctx context.Context, accessToken string) (*jwt.StandardClaims, error)
}

var _ usecase.AuthUseCase = (*AuthUseCase)(nil)

// Login executes OnLogin.
func (mAuthUC AuthUseCase) Login(ctx context.Context, loginInput usecase.AuthLoginInput) (*usecase.AuthTokenOutput, error) {
	return mAuthUC.OnLogin(ctx, loginInput)
}

// Authorize executes Authorize.
func (mAuthUC AuthUseCase) Authorize(ctx context.Context, accessToken string) (*jwt.StandardClaims, error) {
	return mAuthUC.OnAuthorize(ctx, accessToken)
}

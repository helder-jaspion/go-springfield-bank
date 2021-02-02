package mock

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
)

// AuthUseCase mocks an usecase.AuthUseCase.
type AuthUseCase struct {
	OnLogin func(ctx context.Context, loginInput usecase.AuthLoginInput) (*usecase.AuthTokenOutput, error)
}

var _ usecase.AuthUseCase = (*AuthUseCase)(nil)

// Login returns the result of OnLogin.
func (mAuthUC AuthUseCase) Login(ctx context.Context, loginInput usecase.AuthLoginInput) (*usecase.AuthTokenOutput, error) {
	return mAuthUC.OnLogin(ctx, loginInput)
}

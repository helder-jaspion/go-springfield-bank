package mock

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
)

// AccountUseCase mocks an AccountUseCase.
type AccountUseCase struct {
	OnCreate     func(ctx context.Context, accountInput usecase.AccountCreateInput) (*usecase.AccountCreateOutput, error)
	OnFetch      func(ctx context.Context) ([]usecase.AccountFetchOutput, error)
	OnGetBalance func(ctx context.Context, id model.AccountID) (*usecase.AccountBalanceOutput, error)
}

var _ usecase.AccountUseCase = (*AccountUseCase)(nil)

// Create returns the result of OnCreate.
func (m AccountUseCase) Create(ctx context.Context, accountInput usecase.AccountCreateInput) (*usecase.AccountCreateOutput, error) {
	return m.OnCreate(ctx, accountInput)
}

// Fetch returns the result of OnFetch.
func (m AccountUseCase) Fetch(ctx context.Context) ([]usecase.AccountFetchOutput, error) {
	return m.OnFetch(ctx)
}

// GetBalance returns the result of OnGetBalance.
func (m AccountUseCase) GetBalance(ctx context.Context, id model.AccountID) (*usecase.AccountBalanceOutput, error) {
	return m.OnGetBalance(ctx, id)
}

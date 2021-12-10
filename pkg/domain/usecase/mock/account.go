package mock

import (
	"context"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
)

// AccountUseCase mocks an usecase.AccountUseCase.
type AccountUseCase struct {
	OnCreate     func(ctx context.Context, accountInput usecase.AccountCreateInput) (*usecase.AccountCreateOutput, error)
	OnFetch      func(ctx context.Context) ([]usecase.AccountFetchOutput, error)
	OnGetBalance func(ctx context.Context, id model.AccountID) (*usecase.AccountBalanceOutput, error)
}

var _ usecase.AccountUseCase = (*AccountUseCase)(nil)

// Create returns the result of OnCreate.
func (mAccUC AccountUseCase) Create(ctx context.Context, accountInput usecase.AccountCreateInput) (*usecase.AccountCreateOutput, error) {
	return mAccUC.OnCreate(ctx, accountInput)
}

// Fetch returns the result of OnFetch.
func (mAccUC AccountUseCase) Fetch(ctx context.Context) ([]usecase.AccountFetchOutput, error) {
	return mAccUC.OnFetch(ctx)
}

// GetBalance returns the result of OnGetBalance.
func (mAccUC AccountUseCase) GetBalance(ctx context.Context, id model.AccountID) (*usecase.AccountBalanceOutput, error) {
	return mAccUC.OnGetBalance(ctx, id)
}

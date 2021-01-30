package mock

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
)

// AccountUseCaseMock mocks an AccountUseCase
type AccountUseCaseMock struct {
	OnCreate func(ctx context.Context, accountInput usecase.AccountCreateInput) (*usecase.AccountCreateOutput, error)
}

var _ usecase.AccountUseCase = (*AccountUseCaseMock)(nil)

// Create returns the result of OnCreate
func (m AccountUseCaseMock) Create(ctx context.Context, accountInput usecase.AccountCreateInput) (*usecase.AccountCreateOutput, error) {
	return m.OnCreate(ctx, accountInput)
}

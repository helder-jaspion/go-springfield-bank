package usecase

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
)

// AccountUseCaseMock mocks an AccountUseCase
type AccountUseCaseMock struct {
	OnCreate func(ctx context.Context, accountInput AccountCreateInput) (*model.Account, error)
}

var _ AccountUseCase = (*AccountUseCaseMock)(nil)

// Create returns the result of OnCreate
func (m AccountUseCaseMock) Create(ctx context.Context, accountInput AccountCreateInput) (*model.Account, error) {
	return m.OnCreate(ctx, accountInput)
}

package usecase

import (
	"context"
)

// AccountUseCaseMock mocks an AccountUseCase
type AccountUseCaseMock struct {
	OnCreate func(ctx context.Context, accountInput AccountCreateInput) (*AccountCreateOutput, error)
}

var _ AccountUseCase = (*AccountUseCaseMock)(nil)

// Create returns the result of OnCreate
func (m AccountUseCaseMock) Create(ctx context.Context, accountInput AccountCreateInput) (*AccountCreateOutput, error) {
	return m.OnCreate(ctx, accountInput)
}

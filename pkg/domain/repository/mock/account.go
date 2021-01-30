package mock

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

// AccountRepositoryMock mocks an AccountRepository
type AccountRepositoryMock struct {
	OnCreate func(ctx context.Context, account *model.Account) error
}

var _ repository.AccountRepository = (*AccountRepositoryMock)(nil)

// Create returns the result of OnCreate
func (a AccountRepositoryMock) Create(ctx context.Context, account *model.Account) error {
	return a.OnCreate(ctx, account)
}

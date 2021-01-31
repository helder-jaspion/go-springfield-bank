package mock

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

// AccountRepository mocks an AccountRepository
type AccountRepository struct {
	OnCreate      func(ctx context.Context, account *model.Account) error
	OnExistsByCPF func(ctx context.Context, cpf model.CPF) (bool, error)
}

var _ repository.AccountRepository = (*AccountRepository)(nil)

// Create returns the result of OnCreate
func (a AccountRepository) Create(ctx context.Context, account *model.Account) error {
	return a.OnCreate(ctx, account)
}

// ExistsByCPF returns the result of OnExistsByCPF
func (a AccountRepository) ExistsByCPF(ctx context.Context, cpf model.CPF) (bool, error) {
	return a.OnExistsByCPF(ctx, cpf)
}

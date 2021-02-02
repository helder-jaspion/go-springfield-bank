package mock

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

// AccountRepository mocks an AccountRepository.
type AccountRepository struct {
	OnCreate      func(ctx context.Context, account *model.Account) error
	OnExistsByCPF func(ctx context.Context, cpf model.CPF) (bool, error)
	OnGetByCPF    func(ctx context.Context, cpf model.CPF) (*model.Account, error)
	OnFetch       func(ctx context.Context) ([]model.Account, error)
	OnGetBalance  func(ctx context.Context, id model.AccountID) (*model.Account, error)
}

var _ repository.AccountRepository = (*AccountRepository)(nil)

// Create returns the result of OnCreate.
func (mAccRepo AccountRepository) Create(ctx context.Context, account *model.Account) error {
	return mAccRepo.OnCreate(ctx, account)
}

// ExistsByCPF returns the result of OnExistsByCPF.
func (mAccRepo AccountRepository) ExistsByCPF(ctx context.Context, cpf model.CPF) (bool, error) {
	return mAccRepo.OnExistsByCPF(ctx, cpf)
}

// GetByCPF returns the result of OnGetByCPF.
func (mAccRepo AccountRepository) GetByCPF(ctx context.Context, cpf model.CPF) (*model.Account, error) {
	return mAccRepo.OnGetByCPF(ctx, cpf)
}

// Fetch returns the result of OnFetch.
func (mAccRepo AccountRepository) Fetch(ctx context.Context) ([]model.Account, error) {
	return mAccRepo.OnFetch(ctx)
}

// GetBalance returns the result of OnGetBalance.
func (mAccRepo AccountRepository) GetBalance(ctx context.Context, id model.AccountID) (*model.Account, error) {
	return mAccRepo.OnGetBalance(ctx, id)
}

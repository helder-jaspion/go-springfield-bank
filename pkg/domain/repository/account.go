package repository

import (
	"context"
	"errors"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
)

var (
	// ErrAccountNotFound happens when the account was not found based on search params.
	ErrAccountNotFound = errors.New("account not found")
)

// AccountRepository is the interface that wraps account datasource methods.
type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	ExistsByCPF(ctx context.Context, cpf model.CPF) (bool, error)
	GetByCPF(ctx context.Context, cpf model.CPF) (*model.Account, error)
	Fetch(ctx context.Context) ([]model.Account, error)
	GetBalance(ctx context.Context, id model.AccountID) (*model.Account, error)
	UpdateBalance(ctx context.Context, id model.AccountID, balance model.Money) error
}

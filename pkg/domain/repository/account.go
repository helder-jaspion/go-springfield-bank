package repository

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
)

// AccountRepository is the interface that wraps account datasource methods.
type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	ExistsByCPF(ctx context.Context, cpf model.CPF) (bool, error)
}

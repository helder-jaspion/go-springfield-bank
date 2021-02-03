package repository

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
)

// TransferRepository is the interface that wraps transfer datasource methods.
type TransferRepository interface {
	Transaction
	Create(ctx context.Context, transfer *model.Transfer) error
	Fetch(ctx context.Context, accountID model.AccountID) ([]model.Transfer, error)
}

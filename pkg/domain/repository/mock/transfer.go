package mock

import (
	"context"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

// TransferRepository mocks an TransferRepository.
type TransferRepository struct {
	OnCreate            func(ctx context.Context, transfer *model.Transfer) error
	OnFetch             func(ctx context.Context, accountID model.AccountID) ([]model.Transfer, error)
	OnWithinTransaction func(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error)
}

var _ repository.TransferRepository = (*TransferRepository)(nil)

// Create executes OnCreate.
func (mTrfRepo TransferRepository) Create(ctx context.Context, transfer *model.Transfer) error {
	return mTrfRepo.OnCreate(ctx, transfer)
}

// Fetch executes OnFetch.
func (mTrfRepo TransferRepository) Fetch(ctx context.Context, accountID model.AccountID) ([]model.Transfer, error) {
	return mTrfRepo.OnFetch(ctx, accountID)
}

// WithinTransaction executes OnWithinTransaction.
func (mTrfRepo TransferRepository) WithinTransaction(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error) {
	return mTrfRepo.OnWithinTransaction(ctx, txFunc)
}

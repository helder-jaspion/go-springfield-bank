package mock

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
)

// TransferUseCase mocks an usecase.TransferUseCase.
type TransferUseCase struct {
	OnCreate func(ctx context.Context, transferInput usecase.TransferCreateInput) (*usecase.TransferCreateOutput, error)
	OnFetch  func(ctx context.Context) ([]usecase.TransferFetchOutput, error)
}

var _ usecase.TransferUseCase = (*TransferUseCase)(nil)

// Create returns the result of OnCreate.
func (mTrfUC TransferUseCase) Create(ctx context.Context, transferInput usecase.TransferCreateInput) (*usecase.TransferCreateOutput, error) {
	return mTrfUC.OnCreate(ctx, transferInput)
}

// Fetch returns the result of OnFetch.
func (mTrfUC TransferUseCase) Fetch(ctx context.Context, accountID model.AccountID) ([]usecase.TransferFetchOutput, error) {
	return mTrfUC.OnFetch(ctx)
}

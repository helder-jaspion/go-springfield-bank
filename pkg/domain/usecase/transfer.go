package usecase

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

// TransferUseCase is the interface that wraps all business logic methods related to transfers.
type TransferUseCase interface {
	Create(ctx context.Context, transferInput TransferCreateInput) (*TransferCreateOutput, error)
	Fetch(ctx context.Context, accountID model.AccountID) ([]TransferFetchOutput, error)
}

type transferUseCase struct {
	trfRepo repository.TransferRepository
	accRepo repository.AccountRepository
}

// NewTransferUseCase instantiates a new TransferUseCase.
func NewTransferUseCase(trfRepo repository.TransferRepository, accRepo repository.AccountRepository) TransferUseCase {
	return &transferUseCase{
		trfRepo: trfRepo,
		accRepo: accRepo,
	}
}

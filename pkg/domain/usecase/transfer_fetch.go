package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
)

// ErrTransferFetch happens when an error occurred while fetching the transfers.
var ErrTransferFetch = errors.New("could not fetch transfers")

// TransferFetchOutput represents the output data of the fetch method.
type TransferFetchOutput struct {
	TransferCreateOutput
}

func newTransferFetchOutputList(transfers []model.Transfer) []TransferFetchOutput {
	outputs := make([]TransferFetchOutput, 0)

	for _, transfer := range transfers {
		outputs = append(outputs, TransferFetchOutput{
			TransferCreateOutput: *newTransferCreateOutput(&transfer),
		})
	}

	return outputs
}

// Fetch returns all the transfers from repository.TransferRepository.
func (trfUC *transferUseCase) Fetch(ctx context.Context, accountID model.AccountID) ([]TransferFetchOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	transfers, err := trfUC.trfRepo.Fetch(ctx, accountID)
	if err != nil {
		return nil, ErrTransferFetch
	}

	return newTransferFetchOutputList(transfers), nil
}

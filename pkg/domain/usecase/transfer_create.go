package usecase

import (
	"context"
	"errors"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

var (
	// ErrTransferOriginAccountRequired happens when the Transfer origin account ID is not between 2-100 chars long.
	ErrTransferOriginAccountRequired = errors.New("'account_origin_id' is required")
	// ErrTransferDestinationAccountRequired happens when the Transfer destination account ID is not between 2-100 chars long.
	ErrTransferDestinationAccountRequired = errors.New("'account_destination_id' is required")
	// ErrTransferAmountNotPositive happens when the Transfer amount is less or equal to zero.
	ErrTransferAmountNotPositive = errors.New("'amount' must be greater than zero")
	// ErrTransferSameAccount happens when the origin and destination account IDs are the same.
	ErrTransferSameAccount = errors.New("origin and destination accounts must not be the same")
	// ErrAccountCurrentBalanceInsufficient happens when the origin account balance is less than the transfer amount.
	ErrAccountCurrentBalanceInsufficient = errors.New("current account balance is insufficient")
	// ErrTransferCreate happens when an error occurred and the transfer was not created.
	ErrTransferCreate = errors.New("could not create transfer")
)

// TransferCreateInput represents the expected input data when creating a transfer.
type TransferCreateInput struct {
	AccountOriginID      string  `json:"-"`
	AccountDestinationID string  `json:"account_destination_id"`
	Amount               float64 `json:"amount"`
}

// Validate validates the TransferCreateInput fields.
func (input *TransferCreateInput) Validate() error {
	input.AccountOriginID = strings.TrimSpace(input.AccountOriginID)
	if len(input.AccountOriginID) < 1 {
		return ErrTransferOriginAccountRequired
	}

	input.AccountDestinationID = strings.TrimSpace(input.AccountDestinationID)
	if len(input.AccountDestinationID) < 1 {
		return ErrTransferDestinationAccountRequired
	}

	if input.Amount <= 0 {
		return ErrTransferAmountNotPositive
	}

	if input.AccountOriginID == input.AccountDestinationID {
		return ErrTransferSameAccount
	}

	return nil
}

// TransferCreateOutput represents the output data of the create method.
type TransferCreateOutput struct {
	ID                   string    `json:"id"`
	AccountOriginID      string    `json:"account_origin_id"`
	AccountDestinationID string    `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

func newTransferCreateOutput(transfer *model.Transfer) *TransferCreateOutput {
	return &TransferCreateOutput{
		ID:                   string(transfer.ID),
		AccountOriginID:      string(transfer.AccountOriginID),
		AccountDestinationID: string(transfer.AccountDestinationID),
		Amount:               transfer.Amount.Float64(),
		CreatedAt:            transfer.CreatedAt,
	}
}

// Create validates the input, debit the amount from origin account, credit on destination account and saves the transfer.
func (trfUC transferUseCase) Create(ctx context.Context, transferInput TransferCreateInput) (*TransferCreateOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := transferInput.Validate()
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Interface("input", transferInput).Msg("transfer create input is not valid")
		return nil, err
	}

	transfer := model.NewTransfer(
		transferInput.AccountOriginID,
		transferInput.AccountDestinationID,
		transferInput.Amount)

	_, err = trfUC.trfRepo.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		err := trfUC.debitOriginAccount(txCtx, transfer)
		if err != nil {
			return nil, err
		}

		err = trfUC.creditDestinationAccount(txCtx, transfer)
		if err != nil {
			return nil, err
		}

		err = trfUC.trfRepo.Create(txCtx, transfer)
		return nil, err
	})
	if err != nil {
		if err == repository.ErrAccountNotFound || err == ErrAccountCurrentBalanceInsufficient {
			return nil, err
		}
		log.Ctx(ctx).Error().Err(err).Interface("transfer", transfer).Msg("error persisting new transfer")
		return nil, ErrTransferCreate
	}

	return newTransferCreateOutput(transfer), nil
}

func (trfUC transferUseCase) debitOriginAccount(ctx context.Context, transfer *model.Transfer) error {
	originAccount, err := trfUC.accRepo.GetBalance(ctx, transfer.AccountOriginID)
	if err != nil {
		return err
	}
	if originAccount.Balance-transfer.Amount < 0 {
		return ErrAccountCurrentBalanceInsufficient
	}

	return trfUC.accRepo.UpdateBalance(ctx, transfer.AccountOriginID, originAccount.Balance-transfer.Amount)
}

func (trfUC transferUseCase) creditDestinationAccount(ctx context.Context, transfer *model.Transfer) error {
	destinationAccount, err := trfUC.accRepo.GetBalance(ctx, transfer.AccountDestinationID)
	if err != nil {
		return err
	}

	return trfUC.accRepo.UpdateBalance(ctx, transfer.AccountDestinationID, destinationAccount.Balance+transfer.Amount)
}

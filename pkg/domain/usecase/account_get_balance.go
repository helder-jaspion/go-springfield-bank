package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

var (
	// ErrAccountGetBalance happens when an error occurred while getting the account balance.
	ErrAccountGetBalance = errors.New("could not get account balance")
)

// AccountBalanceOutput represents the output data of the GetBalance method.
type AccountBalanceOutput struct {
	ID      string  `json:"id" example:"16b1d860-43d3-4970-bb54-ec395908599a"`
	Balance float64 `json:"balance" example:"9999.99"`
}

func newAccountBalanceOutput(account *model.Account) *AccountBalanceOutput {
	return &AccountBalanceOutput{
		ID:      string(account.ID),
		Balance: account.Balance.Float64(),
	}
}

// GetBalance returns the AccountBalanceOutput with matching ID from repository.AccountRepository.
func (accUC *accountUseCase) GetBalance(ctx context.Context, id model.AccountID) (*AccountBalanceOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	account, err := accUC.accRepo.GetBalance(ctx, id)
	if err != nil {
		if err == repository.ErrAccountNotFound {
			return nil, err
		}
		return nil, ErrAccountGetBalance
	}

	return newAccountBalanceOutput(account), nil
}

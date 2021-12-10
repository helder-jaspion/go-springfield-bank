package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
)

var (
	// ErrAccountFetch happens when an error occurred while fetching the accounts.
	ErrAccountFetch = errors.New("could not fetch accounts")
)

// AccountFetchOutput represents the output data of the fetch method.
type AccountFetchOutput struct {
	AccountCreateOutput
}

func newAccountFetchOutputList(accounts []model.Account) []AccountFetchOutput {
	var outputs = make([]AccountFetchOutput, 0)

	for _, account := range accounts {
		outputs = append(outputs, AccountFetchOutput{
			AccountCreateOutput: *newAccountCreateOutput(&account),
		})
	}

	return outputs
}

// Fetch returns all the accounts from repository.AccountRepository.
func (accUC *accountUseCase) Fetch(ctx context.Context) ([]AccountFetchOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	accounts, err := accUC.accRepo.Fetch(ctx)
	if err != nil {
		return nil, ErrAccountFetch
	}

	return newAccountFetchOutputList(accounts), nil
}

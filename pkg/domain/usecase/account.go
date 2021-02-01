package usecase

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

// AccountUseCase is the interface that wraps all business logic methods related to accounts.
type AccountUseCase interface {
	Create(ctx context.Context, accountInput AccountCreateInput) (*AccountCreateOutput, error)
	Fetch(ctx context.Context) ([]AccountFetchOutput, error)
	GetBalance(ctx context.Context, id model.AccountID) (*AccountBalanceOutput, error)
}

type accountUseCase struct {
	accountRepo repository.AccountRepository
}

// NewAccountUseCase instantiates a new AccountUseCase.
func NewAccountUseCase(accountRepo repository.AccountRepository) AccountUseCase {
	return &accountUseCase{
		accountRepo: accountRepo,
	}
}

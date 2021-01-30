package usecase

import (
	"context"
	"errors"
	"github.com/helder-jaspion/go-springfield-bank/pkg/cpfutil"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"strings"
	"time"
)

// AccountCreateInput represents the expected input data when creating an account.
type AccountCreateInput struct {
	Name    string  `json:"name"`
	CPF     string  `json:"cpf"`
	Secret  string  `json:"secret"`
	Balance float64 `json:"balance"`
}

// AccountUseCase is the interface that wraps all business logic methods related to accounts.
type AccountUseCase interface {
	Create(ctx context.Context, accountInput AccountCreateInput) (*model.Account, error)
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

// Create receives an AccountCreateInput, validates and save it sending to the repository.AccountRepository.
func (accountUC accountUseCase) Create(ctx context.Context, accountInput AccountCreateInput) (*model.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	accountInput.Name = strings.TrimSpace(accountInput.Name)
	if nameLen := len(accountInput.Name); nameLen < 2 || nameLen > 100 {
		return nil, errors.New("'name' must be between 2 and 100 characters in length")
	}

	if !cpfutil.IsValid(accountInput.CPF) {
		return nil, errors.New("'cpf' is invalid")
	}

	if secretLen := len(accountInput.Secret); secretLen < 6 || secretLen > 100 {
		return nil, errors.New("'secret' must be between 6 and 100 characters in length")
	}

	if accountInput.Balance < 0 {
		return nil, errors.New("'balance' must be greater than or equal to zero")
	}

	account := model.NewAccount(accountInput.Name, accountInput.CPF, accountInput.Secret, accountInput.Balance)

	err := account.HashSecret()
	if err != nil {
		return nil, err
	}

	err = accountUC.accountRepo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

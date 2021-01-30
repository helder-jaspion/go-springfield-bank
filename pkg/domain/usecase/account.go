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

var (
	// ErrAccountNameWrongLength happens when the Account Name is not between 2-100 chars long.
	ErrAccountNameWrongLength = errors.New("'name' must be between 2 and 100 characters in length")
	// ErrAccountSecretWrongLength happens when the Account Secret is not between 6-100 chars long.
	ErrAccountSecretWrongLength = errors.New("'secret' must be between 6 and 100 characters in length")
	// ErrAccountBalanceNegative happens when the Account Balance is less than zero.
	ErrAccountBalanceNegative = errors.New("'balance' must be greater than or equal to zero")
	// ErrAccountCPFInvalid happens when the Account CPF is not valid.
	ErrAccountCPFInvalid = errors.New("'cpf' is invalid")
)

// AccountCreateInput represents the expected input data when creating an account.
type AccountCreateInput struct {
	Name    string  `json:"name"`
	CPF     string  `json:"cpf"`
	Secret  string  `json:"secret"`
	Balance float64 `json:"balance"`
}

// Validate validates the AccountCreateInput fields
func (input *AccountCreateInput) Validate() error {
	input.Name = strings.TrimSpace(input.Name)
	if nameLen := len(input.Name); nameLen < 2 || nameLen > 100 {
		return ErrAccountNameWrongLength
	}

	if !cpfutil.IsValid(input.CPF) {
		return ErrAccountCPFInvalid
	}

	if secretLen := len(input.Secret); secretLen < 6 || secretLen > 100 {
		return ErrAccountSecretWrongLength
	}

	if input.Balance < 0 {
		return ErrAccountBalanceNegative
	}

	return nil
}

// AccountCreateOutput represents the output data of the create method.
type AccountCreateOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// AccountUseCase is the interface that wraps all business logic methods related to accounts.
type AccountUseCase interface {
	Create(ctx context.Context, accountInput AccountCreateInput) (*AccountCreateOutput, error)
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
func (accountUC accountUseCase) Create(ctx context.Context, accountInput AccountCreateInput) (*AccountCreateOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := accountInput.Validate()
	if err != nil {
		return nil, err
	}

	account := model.NewAccount(accountInput.Name, accountInput.CPF, accountInput.Secret, accountInput.Balance)

	err = account.HashSecret()
	if err != nil {
		return nil, err
	}

	err = accountUC.accountRepo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return newAccountCreateOutput(account), nil
}

func newAccountCreateOutput(account *model.Account) *AccountCreateOutput {
	return &AccountCreateOutput{
		ID:        account.ID,
		Name:      account.Name,
		CPF:       cpfutil.Format(account.CPF),
		Balance:   account.Balance.Float64(),
		CreatedAt: account.CreatedAt,
	}
}

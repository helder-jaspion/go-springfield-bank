package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/helder-jaspion/go-springfield-bank/pkg/cpfutil"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
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
	// ErrAccountCPFAlreadyExists happens when one tries to create an account with a CPF that is already in use by another account.
	ErrAccountCPFAlreadyExists = errors.New("an account with this CPF already exists")
	// ErrAccountCreate happens when an error occurred and the account was not created.
	ErrAccountCreate = errors.New("could not create account")
)

// AccountCreateInput represents the expected input data when creating an account.
type AccountCreateInput struct {
	Name    string  `json:"name" example:"Bart Simpson"`
	CPF     string  `json:"cpf" example:"999.999.999-99"`
	Secret  string  `json:"secret" example:"S3cr3t"`
	Balance float64 `json:"balance" example:"9999.99" default:"0"`
}

// Validate validates the AccountCreateInput fields.
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
	ID        string    `json:"id" example:"16b1d860-43d3-4970-bb54-ec395908599a"`
	Name      string    `json:"name" example:"Bart Simpson"`
	CPF       string    `json:"cpf" example:"999.999.999-99"`
	Balance   float64   `json:"balance" example:"9999.99"`
	CreatedAt time.Time `json:"created_at" example:"2020-12-31T23:59:59.999999-03:00"`
}

func newAccountCreateOutput(account *model.Account) *AccountCreateOutput {
	return &AccountCreateOutput{
		ID:        string(account.ID),
		Name:      account.Name,
		CPF:       account.CPF.String(),
		Balance:   account.Balance.Float64(),
		CreatedAt: account.CreatedAt,
	}
}

// Create receives an AccountCreateInput, validates and save it sending to the repository.AccountRepository.
func (accUC accountUseCase) Create(ctx context.Context, accountInput AccountCreateInput) (*AccountCreateOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := accountInput.Validate()
	if err != nil {
		accountInput.Secret = "[MASKED]" // removes secret to prevent logging it
		log.Ctx(ctx).Warn().Err(err).Interface("input", accountInput).Msg("account create input is not valid")
		return nil, err
	}

	account := model.NewAccount(accountInput.Name, accountInput.CPF, accountInput.Secret, accountInput.Balance)

	err = account.HashSecret()
	if err != nil {
		account.Secret = "[MASKED]" // removes secret to prevent logging it
		log.Ctx(ctx).Error().Stack().Err(err).Interface("account", account).Msg("error hashing account secret")
		return nil, ErrAccountCreate
	}

	accountExists, err := accUC.accRepo.ExistsByCPF(ctx, account.CPF)
	if err != nil {
		return nil, ErrAccountCreate
	}
	if accountExists {
		return nil, ErrAccountCPFAlreadyExists
	}

	err = accUC.accRepo.Create(ctx, account)
	if err != nil {
		log.Ctx(ctx).Error().Stack().Err(err).Interface("account", account).Msg("error persisting new account")
		return nil, ErrAccountCreate
	}

	return newAccountCreateOutput(account), nil
}

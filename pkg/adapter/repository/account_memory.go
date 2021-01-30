package repository

import (
	"context"
	"errors"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"sync"
)

// AccountMemoryRepository represents an in-memory database to hold accounts.
// As it keeps the data in-memory, the data is lost when the application is shutdown.
type AccountMemoryRepository struct {
	accounts map[string]model.Account
	lock     *sync.RWMutex
}

// NewAccountMemoryRepository instantiates a new account in-memory repository.
func NewAccountMemoryRepository() repository.AccountRepository {
	return &AccountMemoryRepository{
		accounts: map[string]model.Account{},
		lock:     &sync.RWMutex{},
	}
}

// Create adds the account to the accounts map.
// It returns an error if an account with the same id already exists.
func (repo AccountMemoryRepository) Create(_ context.Context, account *model.Account) error {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	_, ok := repo.accounts[account.ID]
	if ok {
		return errors.New("account id already exists")
	}

	repo.accounts[account.ID] = *account

	return nil
}

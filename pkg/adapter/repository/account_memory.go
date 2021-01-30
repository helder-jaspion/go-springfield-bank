package repository

import (
	"context"
	"errors"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"sync"
)

// AccountMemoryRepository represents an in-memory database to hold accountsMap.
// As it keeps the data in-memory, the data is lost when the application is shutdown.
type AccountMemoryRepository struct {
	accountsMap map[string]model.Account
	lock        *sync.RWMutex
}

var _ repository.AccountRepository = (*AccountMemoryRepository)(nil)

// NewAccountMemoryRepository instantiates a new account in-memory repository.
func NewAccountMemoryRepository(accounts ...model.Account) *AccountMemoryRepository {
	accountsMap := make(map[string]model.Account, len(accounts))

	for _, v := range accounts {
		accountsMap[v.ID] = v
	}

	return &AccountMemoryRepository{
		accountsMap: accountsMap,
		lock:        &sync.RWMutex{},
	}
}

// Create adds the account to the accountsMap map.
// It returns an error if an account with the same id already exists.
func (repo AccountMemoryRepository) Create(_ context.Context, account *model.Account) error {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	_, ok := repo.accountsMap[account.ID]
	if ok {
		return errors.New("account id already exists")
	}

	repo.accountsMap[account.ID] = *account

	return nil
}

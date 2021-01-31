package repository

import (
	"context"
	"errors"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"sync"
)

// AccountMemoryRepository represents an in-memory database to hold accountsByIDMap.
// As it keeps the data in-memory, the data is lost when the application is shutdown.
type AccountMemoryRepository struct {
	accountsByIDMap  map[string]model.Account
	accountsByCPFMap map[string]model.Account
	lock             *sync.RWMutex
}

var _ repository.AccountRepository = (*AccountMemoryRepository)(nil)

// NewAccountMemoryRepository instantiates a new account in-memory repository.
func NewAccountMemoryRepository(accounts ...model.Account) *AccountMemoryRepository {
	accountsByIDMap := make(map[string]model.Account, len(accounts))
	accountsByCPFMap := make(map[string]model.Account, len(accounts))

	for _, v := range accounts {
		accountsByIDMap[v.ID] = v
		accountsByCPFMap[v.CPF] = v
	}

	return &AccountMemoryRepository{
		accountsByIDMap:  accountsByIDMap,
		accountsByCPFMap: accountsByCPFMap,
		lock:             &sync.RWMutex{},
	}
}

// Create adds the account to the accountsByIDMap map.
// It returns an error if an account with the same id already exists.
func (repo AccountMemoryRepository) Create(_ context.Context, account *model.Account) error {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	_, ok := repo.accountsByIDMap[account.ID]
	if ok {
		return errors.New("account id already exists")
	}

	_, ok = repo.accountsByCPFMap[account.CPF]
	if ok {
		return errors.New("account cpf already exists")
	}

	repo.accountsByIDMap[account.ID] = *account

	return nil
}

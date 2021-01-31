package memory

import (
	"context"
	"errors"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"sync"
)

// AccountRepository represents an in-memory database to hold accountsByIDMap.
// As it keeps the data in-memory, the data is lost when the application is shutdown.
type AccountRepository struct {
	accountsByIDMap  map[model.AccountID]model.Account
	accountsByCPFMap map[model.CPF]model.Account
	lock             *sync.RWMutex
}

var _ repository.AccountRepository = (*AccountRepository)(nil)

// NewAccountRepository instantiates a new account in-memory repository.
func NewAccountRepository(accounts ...model.Account) *AccountRepository {
	accountsByIDMap := make(map[model.AccountID]model.Account, len(accounts))
	accountsByCPFMap := make(map[model.CPF]model.Account, len(accounts))

	for _, v := range accounts {
		accountsByIDMap[v.ID] = v
		accountsByCPFMap[v.CPF] = v
	}

	return &AccountRepository{
		accountsByIDMap:  accountsByIDMap,
		accountsByCPFMap: accountsByCPFMap,
		lock:             &sync.RWMutex{},
	}
}

// Create adds the account to the accountsByIDMap map.
// It returns an error if an account with the same id already exists.
func (repo AccountRepository) Create(_ context.Context, account *model.Account) error {
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
	repo.accountsByCPFMap[account.CPF] = *account

	return nil
}

// ExistsByCPF verifies if there is an account with the same cpf.
func (repo AccountRepository) ExistsByCPF(_ context.Context, cpf model.CPF) (bool, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	_, ok := repo.accountsByCPFMap[cpf]
	return ok, nil
}

// Fetch returns all the accounts saved.
func (repo AccountRepository) Fetch(_ context.Context) ([]model.Account, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	values := make([]model.Account, 0, len(repo.accountsByIDMap))

	for _, v := range repo.accountsByIDMap {
		values = append(values, v)
	}

	return values, nil
}

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
func (accRepo AccountRepository) Create(_ context.Context, account *model.Account) error {
	accRepo.lock.RLock()
	defer accRepo.lock.RUnlock()

	_, ok := accRepo.accountsByIDMap[account.ID]
	if ok {
		return errors.New("account id already exists")
	}

	_, ok = accRepo.accountsByCPFMap[account.CPF]
	if ok {
		return errors.New("account cpf already exists")
	}

	accRepo.accountsByIDMap[account.ID] = *account
	accRepo.accountsByCPFMap[account.CPF] = *account

	return nil
}

// ExistsByCPF verifies if there is an account with the same cpf.
func (accRepo AccountRepository) ExistsByCPF(_ context.Context, cpf model.CPF) (bool, error) {
	accRepo.lock.RLock()
	defer accRepo.lock.RUnlock()

	_, ok := accRepo.accountsByCPFMap[cpf]
	return ok, nil
}

// GetByCPF returns the account that corresponds to the CPF.
func (accRepo AccountRepository) GetByCPF(_ context.Context, cpf model.CPF) (*model.Account, error) {
	accRepo.lock.RLock()
	defer accRepo.lock.RUnlock()

	account, ok := accRepo.accountsByCPFMap[cpf]
	if !ok {
		return nil, repository.ErrAccountNotFound
	}

	return &account, nil
}

// Fetch returns all the accounts saved.
func (accRepo AccountRepository) Fetch(_ context.Context) ([]model.Account, error) {
	accRepo.lock.RLock()
	defer accRepo.lock.RUnlock()

	values := make([]model.Account, 0, len(accRepo.accountsByIDMap))

	for _, v := range accRepo.accountsByIDMap {
		values = append(values, v)
	}

	return values, nil
}

// GetBalance returns the ID and balance of the account from db.
func (accRepo AccountRepository) GetBalance(_ context.Context, id model.AccountID) (*model.Account, error) {
	accRepo.lock.RLock()
	defer accRepo.lock.RUnlock()

	account, ok := accRepo.accountsByIDMap[id]
	if !ok {
		return nil, repository.ErrAccountNotFound
	}

	return &model.Account{
		ID:      id,
		Balance: account.Balance,
	}, nil
}

// UpdateBalance updates the account balance with the new value.
func (accRepo AccountRepository) UpdateBalance(_ context.Context, id model.AccountID, balance model.Money) error {
	accRepo.lock.RLock()
	defer accRepo.lock.RUnlock()

	account, ok := accRepo.accountsByIDMap[id]
	if !ok {
		return repository.ErrAccountNotFound
	}

	account.Balance = balance
	accRepo.accountsByIDMap[id] = account
	accRepo.accountsByCPFMap[account.CPF] = account

	return nil
}

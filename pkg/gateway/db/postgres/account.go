package postgres

import (
	"context"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type accountRepository struct {
	db *pgxpool.Pool
}

// NewAccountRepository instantiates a new account postgres repository.
func NewAccountRepository(db *pgxpool.Pool) repository.AccountRepository {
	return &accountRepository{db}
}

func (ar *accountRepository) Create(ctx context.Context, account *model.Account) error {
	var query = `
		INSERT INTO
			accounts (id, name, cpf, secret, balance, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6)
	`

	_, err := ar.db.Exec(
		ctx,
		query,
		string(account.ID),
		account.Name,
		account.CPF,
		account.Secret,
		account.Balance,
		account.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ar *accountRepository) ExistsByCPF(ctx context.Context, cpf model.CPF) (bool, error) {
	var query = `SELECT EXISTS(SELECT id FROM accounts WHERE cpf = $1)`

	accountExists := false
	err := ar.db.QueryRow(ctx, query, cpf).Scan(&accountExists)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return accountExists, err
}

func (ar *accountRepository) Fetch(ctx context.Context) ([]model.Account, error) {
	var query = `
		SELECT
			id, name, cpf, secret, balance, created_at
		FROM accounts
	`

	rows, err := ar.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts = make([]model.Account, 0)
	for rows.Next() {
		var account model.Account
		err := rows.Scan(&account.ID, &account.Name, &account.CPF, &account.Secret, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (ar *accountRepository) GetBalance(ctx context.Context, id model.AccountID) (*model.Account, error) {
	var query = "SELECT balance FROM accounts WHERE id = $1"

	account := new(model.Account)
	account.ID = id

	err := ar.db.QueryRow(ctx, query, string(id)).Scan(&account.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

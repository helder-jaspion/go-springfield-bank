package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
)

type transferRepository struct {
	db *pgxpool.Pool
}

// NewTransferRepository instantiates a new transfer postgres repository.
func NewTransferRepository(db *pgxpool.Pool) repository.TransferRepository {
	return &transferRepository{db}
}

func (trfRepo transferRepository) Create(ctx context.Context, transfer *model.Transfer) error {
	query := `
		INSERT INTO
			transfers (id, account_origin_id, account_destination_id, amount, created_at)
		VALUES
			($1, $2, $3, $4, $5)
	`

	_, err := getConnFromCtx(ctx, trfRepo.db).Exec(
		ctx,
		query,
		string(transfer.ID),
		string(transfer.AccountOriginID),
		string(transfer.AccountDestinationID),
		transfer.Amount,
		transfer.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (trfRepo transferRepository) Fetch(ctx context.Context, accountID model.AccountID) ([]model.Transfer, error) {
	query := `
		SELECT
			id, account_origin_id, account_destination_id, amount, created_at
		FROM transfers
		WHERE account_origin_id = $1 OR account_destination_id = $1
		ORDER BY created_at desc
	`

	rows, err := getConnFromCtx(ctx, trfRepo.db).Query(ctx, query, string(accountID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transfers := make([]model.Transfer, 0)
	for rows.Next() {
		var transfer model.Transfer
		err := rows.Scan(&transfer.ID, &transfer.AccountOriginID, &transfer.AccountDestinationID, &transfer.Amount, &transfer.CreatedAt)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, transfer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transfers, nil
}

func (trfRepo transferRepository) WithinTransaction(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error) {
	return execTransaction(ctx, trfRepo.db, txFunc)
}

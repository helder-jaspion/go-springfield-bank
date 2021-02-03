package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type contextKey int

const (
	transactionContextKey contextKey = iota
)

func execTransaction(ctx context.Context, db *pgxpool.Pool, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		p := recover()
		if p != nil || !errors.Is(err, nil) {
			rbErr := tx.Rollback(ctx)
			if rbErr != nil {
				log.Logger.Error().
					Interface("panic", p).
					AnErr("originalErr", err).
					Err(rbErr).
					Msg("error during transaction rollback")
			} else {
				log.Logger.Warn().
					Interface("panic", p).
					AnErr("originalErr", err).
					Msg("transaction rollback executed")
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL READ COMMITTED")
	if err != nil {
		return nil, err
	}

	ctxTx := context.WithValue(ctx, transactionContextKey, tx)
	data, err = txFunc(ctxTx)
	return data, err
}

func getConnFromCtx(ctx context.Context, db *pgxpool.Pool) *pgxpool.Pool {
	tx, ok := ctx.Value(transactionContextKey).(*pgxpool.Pool)
	if !ok {
		return db
	}

	return tx
}

package repository

import (
	"context"
)

// Transaction is the interface that wraps transaction related methods.
type Transaction interface {
	WithinTransaction(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (data interface{}, err error)
}

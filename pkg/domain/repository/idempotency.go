package repository

import (
	"context"
	"time"
)

// IdempotencyRepository is the interface that wraps idempotency datasource methods.
type IdempotencyRepository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, duration time.Duration) error
}

package redis

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/repository"
	"time"
)

type idempotencyRepository struct {
	client *redis.Client
	prefix string
}

// NewIdempotencyRepository instantiates a new idempotency redis repository.
func NewIdempotencyRepository(client *redis.Client) repository.IdempotencyRepository {
	return &idempotencyRepository{client, "_IDEMPOTENCY_"}
}

func (idpRepo idempotencyRepository) Get(ctx context.Context, key string) ([]byte, error) {
	return idpRepo.client.WithContext(ctx).Get(idpRepo.prefix + key).Bytes()
}

func (idpRepo idempotencyRepository) Set(ctx context.Context, key string, value []byte, duration time.Duration) error {
	return idpRepo.client.WithContext(ctx).Set(idpRepo.prefix+key, value, duration).Err()
}

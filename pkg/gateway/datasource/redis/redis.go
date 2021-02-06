package redis

import (
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

// Connect connects to redis server.
func Connect(url string) (*redis.Client, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	err = client.Ping().Err()
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Connected to Redis on %s", url)

	return client, nil
}

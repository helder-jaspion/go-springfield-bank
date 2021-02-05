package redis

import (
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

// Connect connects to redis server.
func Connect(url string, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	err := client.Ping().Err()
	if err != nil {
		panic(err)
	}

	log.Info().Msgf("Connected to Redis on %s", url)

	return client
}

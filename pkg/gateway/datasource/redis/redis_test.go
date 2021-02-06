package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

var testRedisClient *redis.Client

func TestMain(m *testing.M) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not connect to docker")
	}

	opts := dockertest.RunOptions{
		Repository: "redis",
		Tag:        "6-alpine",
		Entrypoint: []string{"redis-server", "--requirepass", "RedisTest2021!"},
	}

	resource, err := dockerPool.RunWithOptions(&opts)
	if err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not start resource")
	}
	_ = resource.Expire(60) // Tell docker to hard kill the container in 60 seconds
	url := fmt.Sprintf("redis://%s:%s@%s:%s", "", "RedisTest2021!", "localhost", resource.GetPort("6379/tcp"))

	if err = dockerPool.Retry(func() error {
		testRedisClient, err = Connect(url)
		return err
	}); err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not connect to docker")
	}

	defer func() {
		err = testRedisClient.Close()
		if err != nil {
			log.Logger.Error().Stack().Err(err).Msg("Could not close redis test connection")
		}
	}()

	code := m.Run()

	if err := dockerPool.Purge(resource); err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not purge resource")
	}

	os.Exit(code)
}

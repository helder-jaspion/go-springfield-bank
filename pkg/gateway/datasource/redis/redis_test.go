package redis

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-redis/redis"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog/log"
)

var (
	testRedisClient *redis.Client
	testRedisURL    string
)

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
	testRedisURL = fmt.Sprintf("redis://%s:%s@%s:%s", "", "RedisTest2021!", "localhost", resource.GetPort("6379/tcp"))

	if err = dockerPool.Retry(func() error {
		testRedisClient, err = Connect(testRedisURL)
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

func TestConnect(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				url: testRedisURL,
			},
			wantErr: false,
		},
		{
			name: "wrong url format should error",
			args: args{
				url: "http://wrong-url",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Connect(tt.args.url)
			defer func() {
				if got != nil {
					got.Close()
				}
			}()
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if err := got.Ping().Err(); err != nil {
					t.Errorf("Connect() error while ping = %v", err)
					return
				}
			}
		})
	}
}

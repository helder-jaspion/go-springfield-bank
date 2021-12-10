package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog/log"

	"github.com/helder-jaspion/go-springfield-bank/config"
	"github.com/helder-jaspion/go-springfield-bank/pkg/gateway/datasource/postgres"
	redisGateway "github.com/helder-jaspion/go-springfield-bank/pkg/gateway/datasource/redis"
)

const (
	contentType     = "Content-Type"
	jsonContentType = "application/json"
)

var testDbPool *pgxpool.Pool
var testRedisClient *redis.Client

func TestMain(m *testing.M) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not connect to docker")
	}

	postgresResource := getPostgresConn(dockerPool)
	defer func() {
		testDbPool.Close()
	}()

	redisResource := getRedisConn(dockerPool)
	defer func() {
		err = testRedisClient.Close()
		if err != nil {
			log.Logger.Error().Stack().Err(err).Msg("Could not close redis test connection")
		}
	}()

	code := m.Run()

	if err := dockerPool.Purge(postgresResource); err != nil {
		log.Logger.Error().Stack().Err(err).Msg("Could not purge postgres resource")
	}
	if err := dockerPool.Purge(redisResource); err != nil {
		log.Logger.Error().Stack().Err(err).Msg("Could not purge redis resource")
	}

	os.Exit(code)
}

func getPostgresConn(dockerPool *dockertest.Pool) *dockertest.Resource {
	confPostgres := config.ConfPostgres{
		Host:                "localhost",
		DbName:              "postgres_test",
		User:                "postgres_test",
		Password:            "secret",
		SslMode:             "prefer",
		PoolMaxConn:         5,
		PoolMaxConnLifetime: 5 * time.Minute,
		Migrate:             true,
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12-alpine",
		Env: []string{
			"POSTGRES_USER=" + confPostgres.User,
			"POSTGRES_PASSWORD=" + confPostgres.Password,
			"POSTGRES_DB=" + confPostgres.DbName,
		},
	}

	resource, err := dockerPool.RunWithOptions(&opts)
	if err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not start resource")
	}
	_ = resource.Expire(60) // Tell docker to hard kill the container in 60 seconds
	confPostgres.Port = resource.GetPort("5432/tcp")

	if err = dockerPool.Retry(func() error {
		testDbPool, err = postgres.ConnectPool(confPostgres.GetDSN(), confPostgres.Migrate)
		return err
	}); err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not connect to docker")
	}
	return resource
}

func getRedisConn(dockerPool *dockertest.Pool) *dockertest.Resource {
	opts := dockertest.RunOptions{
		Repository: "redis",
		Tag:        "6-alpine",
		Entrypoint: []string{"redis-server", "--requirepass", "RedisTest2021!"},
	}

	resource, err := dockerPool.RunWithOptions(&opts)
	if err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not start resource")
	}
	_ = resource.Expire(600) // Tell docker to hard kill the container in 10 minutes
	url := fmt.Sprintf("redis://%s:%s@%s:%s", "", "RedisTest2021!", "localhost", resource.GetPort("6379/tcp"))

	if err = dockerPool.Retry(func() error {
		testRedisClient, err = redisGateway.Connect(url)
		return err
	}); err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not connect to docker")
	}

	return resource
}

func truncateDatabase(t *testing.T) {
	backgroundCtx := context.Background()

	_, err := testDbPool.Exec(backgroundCtx, "DELETE FROM transfers")
	if err != nil {
		t.Errorf("Error truncating transfers table: %v", err)
	}
	_, err = testDbPool.Exec(backgroundCtx, "DELETE FROM accounts")
	if err != nil {
		t.Errorf("Error truncating accounts table: %v", err)
	}
}

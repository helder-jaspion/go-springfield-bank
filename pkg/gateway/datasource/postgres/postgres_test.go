package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog/log"

	"github.com/helder-jaspion/go-springfield-bank/config"
)

var testDbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	confPostgres := config.ConfPostgres{
		Host:                "localhost",
		DBName:              "postgres_test",
		User:                "postgres_test",
		Password:            "secret",
		SslMode:             "disable",
		PoolMaxConn:         5,
		PoolMaxConnLifetime: 5 * time.Minute,
		Migrate:             true,
	}

	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not connect to docker")
	}

	if err = dockerPool.Client.Ping(); err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("could not connect to docker")
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13-alpine",
		Env: []string{
			"POSTGRES_USER=" + confPostgres.User,
			"POSTGRES_PASSWORD=" + confPostgres.Password,
			"POSTGRES_DB=" + confPostgres.DBName,
		},
	}

	resource, err := dockerPool.RunWithOptions(&opts, func(c *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		c.AutoRemove = true
		c.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not start resource")
	}
	_ = resource.Expire(60) // Tell docker to hard kill the container in 60 seconds
	confPostgres.Port = resource.GetPort("5432/tcp")

	if err = dockerPool.Retry(func() error {
		testDbPool, err = ConnectPool(confPostgres)
		return err
	}); err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not connect to docker")
	}

	defer func() {
		testDbPool.Close()
	}()

	code := m.Run()

	if err := dockerPool.Purge(resource); err != nil {
		log.Logger.Fatal().Stack().Err(err).Msg("Could not purge resource")
	}

	os.Exit(code)
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

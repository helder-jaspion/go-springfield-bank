package postgres

import (
	"context"
	"embed"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/helder-jaspion/go-springfield-bank/config"
)

//go:embed migrations
var migrationsFS embed.FS //nolint:gochecknoglobals

// ConnectPool connects do Postgres and returns a *pgxPool.Pool.
func ConnectPool(conf config.ConfPostgres) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(conf.GetDSN())
	if err != nil {
		return nil, errors.Wrap(err, "error configuring the database")
	}
	pgxConfig.ConnConfig.Logger = zerologadapter.NewLogger(log.Logger)

	dbPool, err := pgxpool.ConnectConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to database")
	}

	if conf.Migrate {
		err = RunMigrations(conf.GetURL())
		if err != nil {
			return nil, errors.Wrap(err, "error migrating postgres database")
		}
	}

	return dbPool, nil
}

// RunMigrations executes the database migrations all the way up.
func RunMigrations(connURL string) error {
	source, err := httpfs.New(http.FS(migrationsFS), "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("httpfs", source, connURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Msg("postgres database migration found no changes")
		} else {
			return err
		}
	}

	log.Info().Msg("postgres database migrated successfully")
	return nil
}

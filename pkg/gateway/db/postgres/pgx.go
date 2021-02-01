package postgres

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	pgxDriver "github.com/golang-migrate/migrate/v4/database/postgres"
	// migrate using sql files
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
)

// ConnectPool connects do Postgres and returns a *pgxPool.Pool.
func ConnectPool(dbURL string, runMigrations bool) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("error configuring the database")
	}
	config.ConnConfig.Logger = zerologadapter.NewLogger(log.Logger)

	dbPool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Unable to connect to database")
	}

	if runMigrations {
		err = RunMigrations(config.ConnConfig)
		if err != nil {
			log.Logger.Fatal().Err(err).Msg("error migrating postgres database")
		}
	}

	return dbPool
}

// RunMigrations executes the database migrations all the way up.
func RunMigrations(connConfig *pgx.ConnConfig) error {
	db := stdlib.OpenDB(*connConfig)
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("error closing db connection from migration")
		}
	}()

	drv, err := pgxDriver.WithInstance(db, &pgxDriver.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations/postgres", "postgres", drv)
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

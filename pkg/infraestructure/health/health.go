package health

import (
	"github.com/heptiolabs/healthcheck"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

// RunHealthServer starts and exposes the health endpoints /live and /ready.
func RunHealthServer(healthPort string, dbPool *pgxpool.Pool) {
	// Create a Handler that we can use to register liveness and readiness checks.
	health := healthcheck.NewHandler()

	db := stdlib.OpenDB(*dbPool.Config().ConnConfig)

	dbCheck := healthcheck.DatabasePingCheck(db, 1*time.Second)

	health.AddReadinessCheck("database", dbCheck)
	health.AddLivenessCheck("database", dbCheck)

	err := http.ListenAndServe(":"+healthPort, health)
	if err != nil {
		log.Error().Err(err).Msg("Could not start health server")
	}
}

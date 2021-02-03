package monitoring

import (
	"github.com/heptiolabs/healthcheck"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

// RunServer starts and exposes the monitoring (health and metrics) endpoints.
func RunServer(port string, dbPool *pgxpool.Pool) {
	adminMux := http.NewServeMux()

	adminMux.Handle("/metrics", promhttp.Handler())

	health := healthHandler(dbPool)
	adminMux.HandleFunc("/live", health.LiveEndpoint)
	adminMux.HandleFunc("/ready", health.ReadyEndpoint)

	err := http.ListenAndServe(":"+port, adminMux)
	if err != nil {
		log.Error().Err(err).Msg("Could not start monitoring server")
	}
}

func healthHandler(dbPool *pgxpool.Pool) healthcheck.Handler {
	// Create a Handler that we can use to register liveness and readiness checks.
	health := healthcheck.NewHandler()

	db := stdlib.OpenDB(*dbPool.Config().ConnConfig)

	dbCheck := healthcheck.DatabasePingCheck(db, 1*time.Second)

	health.AddReadinessCheck("database", dbCheck)
	health.AddLivenessCheck("database", dbCheck)
	return health
}

package monitoring

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/heptiolabs/healthcheck"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

// RunServer starts and exposes the monitoring (health and metrics) endpoints.
func RunServer(port string, dbPool *pgxpool.Pool, redisClient *redis.Client) {
	adminMux := http.NewServeMux()

	adminMux.Handle("/metrics", promhttp.Handler())

	health := healthHandler(dbPool, redisClient)
	adminMux.HandleFunc("/live", health.LiveEndpoint)
	adminMux.HandleFunc("/ready", health.ReadyEndpoint)

	err := http.ListenAndServe(":"+port, adminMux)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Could not start monitoring server")
	}
}

func healthHandler(dbPool *pgxpool.Pool, redisClient *redis.Client) healthcheck.Handler {
	// Create a Handler that we can use to register liveness and readiness checks.
	health := healthcheck.NewHandler()

	dbCheck := PgxPoolSelectCheck(dbPool, 1*time.Second)
	health.AddReadinessCheck("database", dbCheck)
	health.AddLivenessCheck("database", dbCheck)

	redisCheck := RedisPingCheck(redisClient, 1*time.Second)
	health.AddReadinessCheck("redis", redisCheck)
	health.AddLivenessCheck("redis", redisCheck)
	return health
}

// PgxPoolSelectCheck returns a Check that validates connectivity to a database/pgxpool.Pool using 'SELECT 1' query.
func PgxPoolSelectCheck(database *pgxpool.Pool, timeout time.Duration) healthcheck.Check {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if database == nil {
			return fmt.Errorf("database is nil")
		}
		var v interface{}
		return database.QueryRow(ctx, "SELECT 1").Scan(v)
	}
}

// RedisPingCheck returns a Check that validates connectivity to a Redis server using Ping().
func RedisPingCheck(client *redis.Client, timeout time.Duration) healthcheck.Check {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if client == nil {
			return fmt.Errorf("client is nil")
		}

		return client.WithContext(ctx).Ping().Err()
	}
}

package http

import (
	"context"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// StartServer runs the server
func StartServer(server *http.Server) {
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	go gracefulShutdown(server, quit, done, 5*time.Second)

	log.Info().Msgf("Server is ready to handle requests at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Stack().Err(err).Msgf("Could not listen on %s", server.Addr)
	}

	<-done
	log.Info().Msg("Server stopped")
}

func gracefulShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool, gracePeriod time.Duration) {
	<-quit
	log.Info().Msg("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Stack().Err(err).Msg("Could not gracefully shutdown the server")
	}
	close(done)
}

func handleOPTIONS(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Access-Control-Request-Method") != "" {
		// Set CORS headers
		header := w.Header()
		header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
		header.Set("Access-Control-Allow-Origin", "*")
	}

	// Adjust status code to 204
	w.WriteHeader(http.StatusNoContent)
}

func handlePanic(w http.ResponseWriter, r *http.Request, p interface{}) {
	hlog.FromRequest(r).Error().Stack().Interface("panic", p).Msg("Panic recovered")
	w.WriteHeader(http.StatusInternalServerError)
}

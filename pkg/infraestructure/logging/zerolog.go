package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	stdlog "log"
	"os"
	"time"
)

// InitZerolog configures zero log with the initial values
func InitZerolog(levelStr, outputType string) {
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		log.Error().Err(err).Msg("could not parse log level")
	}
	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if outputType == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    true,
			TimeFormat: time.RFC3339,
		})
	}

	log.Logger = log.Logger.With().Caller().Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}

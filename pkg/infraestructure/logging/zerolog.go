package logging

import (
	stdLog "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// InitZeroLog configures zero log with the initial values
func InitZeroLog(levelStr, outputType string) {
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not parse log level")
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

	log.Logger = log.Logger.With().CallerWithSkipFrameCount(2).Logger()

	stdLog.SetFlags(0)
	stdLog.SetOutput(log.Logger)
}

package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	stdLog "log"
	"os"
	"time"
)

// InitZeroLog configures zero log with the initial values
func InitZeroLog(levelStr, outputType string) {
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

	log.Logger = log.Logger.With().CallerWithSkipFrameCount(2).Logger()

	stdLog.SetFlags(0)
	stdLog.SetOutput(log.Logger)
}

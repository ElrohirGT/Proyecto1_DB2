package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ConfigureLogger sets up the zerolog output based on the provided configuration
func ConfigureLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// Use a human-readable console output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

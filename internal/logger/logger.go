// Package logger configures the application-wide zerolog logger.
//
// A single logger is created at startup and passed via dependency injection.
// This package provides no global logger variable — every component that needs
// logging receives a *zerolog.Logger from its constructor.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// New creates and returns a configured zerolog.Logger.
//
//   - level:  minimum log level string ("trace", "debug", "info", "warn", "error")
//   - pretty: when true, writes coloured human-readable output (development only)
func New(level string, pretty bool) zerolog.Logger {
	// Parse the log level string; default to Info on unknown values.
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	// Choose writer: human-readable console output or structured JSON.
	var w io.Writer
	if pretty {
		w = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}
	} else {
		w = os.Stderr
	}

	return zerolog.New(w).
		With().
		Timestamp().
		Caller().
		Logger()
}

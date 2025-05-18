// Package logger provides a thread-safe singleton logger instance
// with centralized configuration for the application.
package logger

import (
	"go.uber.org/zap"
)

// Log is the global logger instance implementing the Logger interface.
var Log *zap.Logger = zap.NewNop()

// LogInit configures the global Log instance.
func LogInit(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	if cfg.Level.Level() != zap.DebugLevel {
		cfg.DisableCaller = true
	}

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

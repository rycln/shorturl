package logger

import (
	"go.uber.org/zap"
)

// для логгера приемлимо использовать глобальную переменную?
var Log *zap.Logger = zap.NewNop()

func LogInit() error {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zap.InfoLevel)
	cfg.DisableCaller = true

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	config "github.com/rycln/shorturl/configs"
	"github.com/rycln/shorturl/internal/app/logger"
	"go.uber.org/zap"
)

func main() {
	err := logger.LogInit()
	if err != nil {
		log.Fatalf("Can't initialize the logger: %v", err)
	}
	defer logger.Log.Sync()

	cfg := config.NewCfg()
	logger.Log.Info("Server configuration:",
		zap.String("addr", cfg.GetServerAddr()),
		zap.String("base_url", cfg.GetBaseAddr()),
		zap.String("storage", cfg.StorageIs()),
	)

	app := fiber.New()

	switch cfg.StorageIs() {
	case "db":
		logger.Log.Info("Storage configuration",
			zap.String("db_dsn", cfg.GetDatabaseDsn()),
		)
		startWithDatabaseStorage(app, cfg)
	case "file":
		logger.Log.Info("Storage configuration",
			zap.String("file_path", cfg.GetFilePath()),
		)
		startWithFileStorage(app, cfg)
	default:
		startWithSimpleStorage(app, cfg)
	}
}

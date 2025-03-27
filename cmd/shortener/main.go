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
	hs, closeStrg := newHandlersSet(cfg)
	if closeStrg != nil {
		defer closeStrg()
	}
	routing(app, cfg, hs)

	err = app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

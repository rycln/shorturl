package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	config "github.com/rycln/shorturl/configs"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/server"
)

func main() {
	err := logger.LogInit()
	if err != nil {
		log.Fatalf("Can't initialize the logger: %v", err)
	}
	defer logger.Log.Sync()

	cfg := config.NewCfg()
	app := fiber.New()

	switch cfg.StorageIs() {
	case "db":
		server.StartWithDatabaseStorage(app, cfg)
	case "file":
		server.StartWithFileStorage(app, cfg)
	default:
		server.StartWithSimpleStorage(app, cfg)
	}
}

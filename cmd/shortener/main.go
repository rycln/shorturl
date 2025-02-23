package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	config "github.com/rycln/shorturl/configs"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/server"
	"github.com/rycln/shorturl/internal/app/storage"
)

func main() {
	err := logger.LogInit()
	if err != nil {
		log.Fatalf("Can't initialize the logger: %v", err)
	}
	defer logger.Log.Sync()

	app := fiber.New()
	storage := storage.NewSimpleMemStorage()
	config := config.NewCfg()
	sa := server.NewServerArgs(storage, config)
	server.Set(app, sa)

	err = app.Listen(config.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

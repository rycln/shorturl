package main

import (
	"github.com/gofiber/fiber/v2"
	config "github.com/rycln/shorturl/configs"
	"github.com/rycln/shorturl/internal/app/server"
	"github.com/rycln/shorturl/internal/app/storage"
)

func main() {
	app := fiber.New()
	storage := storage.NewSimpleMemStorage()
	config := config.NewCfg()
	sa := server.NewServerArgs(storage, config)
	server.Set(app, sa)

	err := app.Listen(config.GetServerAddr())
	if err != nil {
		panic(err)
	}
}

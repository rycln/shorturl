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

	cfg := config.NewCfg()
	strg := storage.NewSimpleMemStorage()

	fd, err := storage.NewFileDecoder(cfg.GetFilePath())
	if err != nil {
		log.Fatalf("Can't open the file: %v", err)
	}
	fd.RestoreStorage(strg)
	fd.Close()

	fe, err := storage.NewFileEncoder(cfg.GetFilePath())
	if err != nil {
		log.Fatalf("Can't open the file: %v", err)
	}
	defer fe.Close()

	app := fiber.New()
	sa := server.NewServerArgs(strg, cfg, fe)
	server.Set(app, sa)

	err = app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

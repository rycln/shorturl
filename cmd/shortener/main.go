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
	app := fiber.New()

	switch {
	case cfg.GetDatabaseDsn() != "":
		err := storage.DBInitPostgre(cfg.GetDatabaseDsn())
		if err != nil {
			log.Fatalf("Can't open database: %v", err)
		}
		db := storage.NewDatabaseStorage(storage.DB)
		defer db.Close()

		sa := server.NewServerArgs(db, cfg)
		server.Set(app, sa)

		err = app.Listen(cfg.GetServerAddr())
		if err != nil {
			log.Fatalf("Can't start the server: %v", err)
		}
	case cfg.GetFilePath() != "":
		fd, err := storage.NewFileDecoder(cfg.GetFilePath())
		if err != nil {
			log.Fatalf("Can't open the file: %v", err)
		}
		defer fd.Close()

		fe, err := storage.NewFileEncoder(cfg.GetFilePath())
		if err != nil {
			log.Fatalf("Can't open the file: %v", err)
		}
		defer fe.Close()

		fs := storage.NewFileStorage(fe, fd)

		sa := server.NewServerArgs(fs, cfg)
		server.Set(app, sa)

		err = app.Listen(cfg.GetServerAddr())
		if err != nil {
			log.Fatalf("Can't start the server: %v", err)
		}
	default:
		strg := storage.NewSimpleStorage()

		sa := server.NewServerArgs(strg, cfg)
		server.Set(app, sa)

		err = app.Listen(cfg.GetServerAddr())
		if err != nil {
			log.Fatalf("Can't start the server: %v", err)
		}
	}
}

package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	storage.DB, err = sql.Open("pgx", cfg.GetDatabaseDsn())
	if err != nil {
		log.Fatalf("Can't open database: %v", err)
	}
	defer storage.DB.Close()

	strg := storage.NewSimpleMemStorage()

	fd, err := storage.NewFileDecoder(cfg.GetFilePath())
	if err != nil {
		log.Fatalf("Can't open the file: %v", err)
	}
	err = fd.RestoreStorage(strg)
	if err != nil {
		log.Fatalf("Can't restore from file: %v", err)
	}
	defer fd.Close()

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

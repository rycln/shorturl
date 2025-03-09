package server

import (
	"log"
	"net/http"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	config "github.com/rycln/shorturl/configs"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap/zapcore"
)

func Set(app *fiber.App, sa *ServerArgs) {
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
		Fields: []string{"url", "method", "latency", "status", "bytesSent"},
		Levels: []zapcore.Level{zapcore.InfoLevel},
	}))

	app.Post("/api/shorten/batch", sa.ShortenBatch)
	app.Post("/api/shorten", sa.ShortenAPI)
	app.Get("/ping", sa.PingDB)
	app.Get("/:short", sa.ReturnURL)
	app.Post("/", sa.ShortenURL)

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusBadRequest)
	})
}

func StartWithSimpleStorage(app *fiber.App, cfg *config.Cfg) {
	strg := storage.NewSimpleStorage()

	sa := NewServerArgs(strg, cfg)
	Set(app, sa)

	err := app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

func StartWithFileStorage(app *fiber.App, cfg *config.Cfg) {
	fs, err := storage.NewFileStorage(cfg.GetFilePath())
	if err != nil {
		log.Fatalf("Can't open the file: %v", err)
	}
	defer fs.Close()

	sa := NewServerArgs(fs, cfg)
	Set(app, sa)

	err = app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

func StartWithDatabaseStorage(app *fiber.App, cfg *config.Cfg) {
	err := storage.DBInitPostgre(cfg.GetDatabaseDsn())
	if err != nil {
		log.Fatalf("Can't open database: %v", err)
	}
	db := storage.NewDatabaseStorage(storage.DB)
	defer db.Close()

	sa := NewServerArgs(db, cfg)
	Set(app, sa)

	err = app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

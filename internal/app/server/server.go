package server

import (
	"log"
	"net/http"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	config "github.com/rycln/shorturl/configs"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/myhash"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap/zapcore"
)

func Set(app *fiber.App, sa *ServerArgs) {
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
		Fields: []string{"url", "method", "latency", "status", "bytesSent"},
		Levels: []zapcore.Level{zapcore.InfoLevel},
	}))

	app.Post("/api/shorten/batch", timeout.NewWithContext(sa.ShortenBatch, sa.cfg.TimeoutDuration()))
	app.Post("/api/shorten", timeout.NewWithContext(sa.ShortenAPI, sa.cfg.TimeoutDuration()))
	app.Get("/ping", timeout.NewWithContext(sa.PingDB, sa.cfg.TimeoutDuration()))
	app.Get("/:short", timeout.NewWithContext(sa.ReturnURL, sa.cfg.TimeoutDuration()))
	app.Post("/", timeout.NewWithContext(sa.ShortenURL, sa.cfg.TimeoutDuration()))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusBadRequest)
	})
}

func StartWithSimpleStorage(app *fiber.App, cfg *config.Cfg) {
	strg := storage.NewSimpleStorage()

	sa := NewServerArgs(strg, cfg, myhash.Base62)
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

	sa := NewServerArgs(fs, cfg, myhash.Base62)
	Set(app, sa)

	err = app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

func StartWithDatabaseStorage(app *fiber.App, cfg *config.Cfg) {
	db, err := storage.NewDB(cfg.GetDatabaseDsn())
	if err != nil {
		log.Fatalf("Can't open database: %v", err)
	}
	defer db.Close()

	err = storage.InitDB(db, cfg.TimeoutDuration())
	if err != nil {
		log.Fatalf("Can't init database: %v", err)
	}

	sdb := storage.NewDatabaseStorage(db)

	sa := NewServerArgs(sdb, cfg, myhash.Base62)
	Set(app, sa)

	err = app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

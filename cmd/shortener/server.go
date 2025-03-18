package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	config "github.com/rycln/shorturl/configs"
	"github.com/rycln/shorturl/internal/app/handlers"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/myhash"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap/zapcore"
)

type handler interface {
	Handle(*fiber.Ctx) error
}

func serverInit(app *fiber.App, to time.Duration, shorten, apiShorten, shortenBatch, retrieve, ping handler) {
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
		Fields: []string{"url", "method", "latency", "status", "bytesSent"},
		Levels: []zapcore.Level{zapcore.InfoLevel},
	}))

	app.Post("/api/shorten/batch", timeout.NewWithContext(shortenBatch.Handle, to))
	app.Post("/api/shorten", timeout.NewWithContext(apiShorten.Handle, to))
	app.Get("/ping", timeout.NewWithContext(ping.Handle, to))
	app.Get("/:short", timeout.NewWithContext(retrieve.Handle, to))
	app.Post("/", timeout.NewWithContext(shorten.Handle, to))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusBadRequest)
	})
}

func startWithSimpleStorage(app *fiber.App, cfg *config.Cfg) {
	strg := storage.NewSimpleStorage()

	shorten := handlers.NewShorten(strg, cfg, myhash.Base62)
	apiShorten := handlers.NewAPIShorten(strg, cfg, myhash.Base62)
	shortenBatch := handlers.NewShortenBatch(strg, cfg, myhash.Base62)
	retrieve := handlers.NewRetrieve(strg)
	ping := handlers.NewPing(strg)

	serverInit(app, cfg.TimeoutDuration(), shorten, apiShorten, shortenBatch, retrieve, ping)

	err := app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

func startWithFileStorage(app *fiber.App, cfg *config.Cfg) {
	fs, err := storage.NewFileStorage(cfg.GetFilePath())
	if err != nil {
		log.Fatalf("Can't open the file: %v", err)
	}
	defer fs.Close()

	shorten := handlers.NewShorten(fs, cfg, myhash.Base62)
	apiShorten := handlers.NewAPIShorten(fs, cfg, myhash.Base62)
	shortenBatch := handlers.NewShortenBatch(fs, cfg, myhash.Base62)
	retrieve := handlers.NewRetrieve(fs)
	ping := handlers.NewPing(fs)

	serverInit(app, cfg.TimeoutDuration(), shorten, apiShorten, shortenBatch, retrieve, ping)

	err = app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

func startWithDatabaseStorage(app *fiber.App, cfg *config.Cfg) {
	db, err := storage.NewDB(cfg.GetDatabaseDsn())
	if err != nil {
		log.Fatalf("Can't open database: %v", err)
	}
	defer db.Close()

	err = storage.InitDB(db, cfg.TimeoutDuration())
	if err != nil {
		log.Fatalf("Can't init database: %v", err)
	}

	dbs := storage.NewDatabaseStorage(db)

	shorten := handlers.NewShorten(dbs, cfg, myhash.Base62)
	apiShorten := handlers.NewAPIShorten(dbs, cfg, myhash.Base62)
	shortenBatch := handlers.NewShortenBatch(dbs, cfg, myhash.Base62)
	retrieve := handlers.NewRetrieve(dbs)
	ping := handlers.NewPing(dbs)

	serverInit(app, cfg.TimeoutDuration(), shorten, apiShorten, shortenBatch, retrieve, ping)

	err = app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

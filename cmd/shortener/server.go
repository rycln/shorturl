package main

import (
	"context"
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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type handlersSet struct {
	shorten       func(*fiber.Ctx) error
	retrieve      func(*fiber.Ctx) error
	apiShorten    func(*fiber.Ctx) error
	shortenBatch  func(*fiber.Ctx) error
	ping          func(*fiber.Ctx) error
	retrieveBatch func(*fiber.Ctx) error
	deleteBatch   func(*fiber.Ctx) error
}

type shutdownFunc func() error

func newHandlersSet(cfg *config.Cfg) (*handlersSet, shutdownFunc) {
	switch cfg.StorageIs() {
	case "db":
		logger.Log.Info("Storage configuration",
			zap.String("db_dsn", cfg.GetDatabaseDsn()),
		)
		ctx, cancel := context.WithCancel(context.Background())
		dbs, close := storage.NewDatabaseStorage(cfg.GetDatabaseDsn(), cfg.GetTimeoutDuration())
		shutdown := func() error {
			cancel()
			return close()
		}
		return &handlersSet{
			shorten:       handlers.NewShortenHandler(dbs, cfg, myhash.Base62),
			retrieve:      handlers.NewRetrieveHandler(dbs),
			apiShorten:    handlers.NewAPIShortenHandler(dbs, cfg, myhash.Base62),
			shortenBatch:  handlers.NewShortenBatchHandler(dbs, cfg, myhash.Base62),
			ping:          handlers.NewPingHandler(dbs),
			retrieveBatch: handlers.NewRetrieveBatchHandler(dbs, cfg),
			deleteBatch:   handlers.NewDeleteBatchHandler(ctx, dbs, cfg),
		}, shutdown
	case "file":
		logger.Log.Info("Storage configuration",
			zap.String("file_path", cfg.GetFilePath()),
		)
		fs, close := storage.NewFileStorage(cfg.GetFilePath())
		return &handlersSet{
			shorten:       handlers.NewShortenHandler(fs, cfg, myhash.Base62),
			retrieve:      handlers.NewRetrieveHandler(fs),
			apiShorten:    handlers.NewAPIShortenHandler(fs, cfg, myhash.Base62),
			shortenBatch:  handlers.NewShortenBatchHandler(fs, cfg, myhash.Base62),
			retrieveBatch: handlers.NewRetrieveBatchHandler(fs, cfg),
		}, close
	default:
		strg := storage.NewSimpleStorage()
		return &handlersSet{
			shorten:      handlers.NewShortenHandler(strg, cfg, myhash.Base62),
			retrieve:     handlers.NewRetrieveHandler(strg),
			apiShorten:   handlers.NewAPIShortenHandler(strg, cfg, myhash.Base62),
			shortenBatch: handlers.NewShortenBatchHandler(strg, cfg, myhash.Base62),
		}, nil
	}
}

func routing(app *fiber.App, hs *handlersSet, to time.Duration) {
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
		Fields: []string{"url", "method", "latency", "status", "bytesSent"},
		Levels: []zapcore.Level{zapcore.InfoLevel},
	}))

	if hs.shortenBatch != nil {
		app.Post("/api/shorten/batch", timeout.NewWithContext(hs.shortenBatch, to))
	}
	if hs.apiShorten != nil {
		app.Post("/api/shorten", timeout.NewWithContext(hs.apiShorten, to))
	}
	if hs.retrieveBatch != nil {
		app.Get("/api/user/urls", timeout.NewWithContext(hs.retrieveBatch, to))
	}
	if hs.deleteBatch != nil {
		app.Delete("/api/user/urls", timeout.NewWithContext(hs.deleteBatch, to))
	}
	if hs.ping != nil {
		app.Get("/ping", timeout.NewWithContext(hs.ping, to))
	}
	if hs.retrieve != nil {
		app.Get("/:short", timeout.NewWithContext(hs.retrieve, to))
	}
	if hs.shorten != nil {
		app.Post("/", timeout.NewWithContext(hs.shorten, to))
	}

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusBadRequest)
	})
}

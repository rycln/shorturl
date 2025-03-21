package main

import (
	"log"
	"net/http"

	"github.com/gofiber/contrib/fiberzap/v2"
	jwtware "github.com/gofiber/contrib/jwt"
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

func startWithSimpleStorage(app *fiber.App, cfg *config.Cfg) {
	strg := storage.NewSimpleStorage()

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
		Fields: []string{"url", "method", "latency", "status", "bytesSent"},
		Levels: []zapcore.Level{zapcore.InfoLevel},
	}))

	to := cfg.TimeoutDuration()
	app.Post("/api/shorten/batch", timeout.NewWithContext(handlers.NewShortenBatch(strg, cfg, myhash.Base62).Handle, to))
	app.Post("/api/shorten", timeout.NewWithContext(handlers.NewAPIShorten(strg, cfg, myhash.Base62).Handle, to))
	app.Get("/:short", timeout.NewWithContext(handlers.NewRetrieve(strg).Handle, to))
	app.Post("/", timeout.NewWithContext(handlers.NewShorten(strg, cfg, myhash.Base62).Handle, to))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusBadRequest)
	})

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

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
		Fields: []string{"url", "method", "latency", "status", "bytesSent"},
		Levels: []zapcore.Level{zapcore.InfoLevel},
	}))

	to := cfg.TimeoutDuration()
	app.Post("/api/shorten/batch", timeout.NewWithContext(handlers.NewShortenBatch(fs, cfg, myhash.Base62).Handle, to))
	app.Post("/api/shorten", timeout.NewWithContext(handlers.NewAPIShorten(fs, cfg, myhash.Base62).Handle, to))
	app.Get("/:short", timeout.NewWithContext(handlers.NewRetrieve(fs).Handle, to))
	app.Post("/", timeout.NewWithContext(handlers.NewShorten(fs, cfg, myhash.Base62).Handle, to))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusBadRequest)
	})

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

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
		Fields: []string{"url", "method", "latency", "status", "bytesSent"},
		Levels: []zapcore.Level{zapcore.InfoLevel},
	}))

	to := cfg.TimeoutDuration()
	app.Post("/api/shorten/batch", timeout.NewWithContext(handlers.NewShortenBatch(dbs, cfg, myhash.Base62).Handle, to))
	app.Post("/api/shorten", timeout.NewWithContext(handlers.NewAPIShorten(dbs, cfg, myhash.Base62).Handle, to))
	app.Get("/ping", timeout.NewWithContext(handlers.NewPing(dbs).Handle, to))
	app.Get("/:short", timeout.NewWithContext(handlers.NewRetrieve(dbs).Handle, to))
	app.Post("/", timeout.NewWithContext(handlers.NewShorten(dbs, cfg, myhash.Base62).Handle, to))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.GetKey())},
	}))

	app.Get("/api/user/urls", timeout.NewWithContext(handlers.NewRetrieveBatch(dbs, cfg).Handle, to))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusBadRequest)
	})

	err = app.Listen(cfg.GetServerAddr())
	if err != nil {
		log.Fatalf("Can't start the server: %v", err)
	}
}

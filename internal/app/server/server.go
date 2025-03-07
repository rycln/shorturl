package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/myhash"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap/zapcore"
)

type configer interface {
	GetBaseAddr() string
	GetDatabaseDsn() string
}

type storager interface {
	AddURL(context.Context, ...storage.ShortenedURL) error
	GetURL(context.Context, string) (string, error)
}

type ServerArgs struct {
	strg storager
	cfg  configer
}

func NewServerArgs(strg storager, cfg configer) *ServerArgs {
	return &ServerArgs{
		strg: strg,
		cfg:  cfg,
	}
}

func (sa *ServerArgs) ShortenURL(c *fiber.Ctx) error {
	body := string(c.Body())
	_, err := url.ParseRequestURI(body)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	origURL := body
	shortURL := myhash.Base62(origURL)

	ctx, cancel := context.WithTimeout(c.Context(), 1*time.Second)
	defer cancel()

	surl := storage.NewShortenedURL(shortURL, origURL)
	err = sa.strg.AddURL(ctx, surl)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	c.Set("Content-Type", "text/plain")
	baseAddr := sa.cfg.GetBaseAddr()
	return c.Status(http.StatusCreated).SendString(baseAddr + "/" + shortURL)
}

func (sa *ServerArgs) ReturnURL(c *fiber.Ctx) error {
	shortURL := c.Params("short")

	ctx, cancel := context.WithTimeout(c.Context(), 1*time.Second)
	defer cancel()

	origURL, err := sa.strg.GetURL(ctx, shortURL)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	c.Set("Location", origURL)
	return c.SendStatus(http.StatusTemporaryRedirect)
}

type apiReq struct {
	URL string `json:"url"`
}

type apiRes struct {
	Result string `json:"result"`
}

func (sa *ServerArgs) ShortenAPI(c *fiber.Ctx) error {
	if !c.Is("json") {
		return c.SendStatus(http.StatusBadRequest)
	}

	var req apiReq
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	_, err = url.ParseRequestURI(req.URL)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	origURL := req.URL
	shortURL := myhash.Base62(origURL)

	ctx, cancel := context.WithTimeout(c.Context(), 1*time.Second)
	defer cancel()

	surl := storage.NewShortenedURL(shortURL, origURL)
	err = sa.strg.AddURL(ctx, surl)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	var res apiRes
	baseAddr := sa.cfg.GetBaseAddr()
	res.Result = baseAddr + "/" + shortURL
	resBody, err := json.Marshal(&res)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
	}
	c.Set("Content-Type", "application/json")
	return c.Status(http.StatusCreated).Send(resBody)
}

func (sa *ServerArgs) PingDB(c *fiber.Ctx) error {
	if sa.cfg.GetDatabaseDsn() == "" {
		return c.SendStatus(http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(c.Context(), 1*time.Second)
	defer cancel()
	if err := storage.DB.PingContext(ctx); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}

func Set(app *fiber.App, sa *ServerArgs) {
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
		Fields: []string{"url", "method", "latency", "status", "bytesSent"},
		Levels: []zapcore.Level{zapcore.InfoLevel},
	}))

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

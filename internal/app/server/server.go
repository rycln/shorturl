package server

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/myhash"
	"go.uber.org/zap/zapcore"
)

type Configer interface {
	GetBaseAddr() string
}

type Storager interface {
	AddURL(string, string)
	GetURL(string) (string, error)
}

type ServerArgs struct {
	storage Storager
	config  Configer
}

func NewServerArgs(storage Storager, config Configer) *ServerArgs {
	return &ServerArgs{
		storage: storage,
		config:  config,
	}
}

func (sa *ServerArgs) ShortenURL(c *fiber.Ctx) error {
	body := string(c.Body())
	_, err := url.ParseRequestURI(body)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	fullURL := body
	shortURL := myhash.Base62(fullURL)
	sa.storage.AddURL(shortURL, fullURL)
	c.Set("Content-Type", "text/plain")
	baseAddr := sa.config.GetBaseAddr()
	return c.Status(http.StatusCreated).SendString(baseAddr + "/" + shortURL)
}

func (sa *ServerArgs) ReturnURL(c *fiber.Ctx) error {
	shortURL := c.Params("short")
	fullURL, err := sa.storage.GetURL(shortURL)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	c.Set("Location", fullURL)
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

	fullURL := req.URL
	shortURL := myhash.Base62(fullURL)
	sa.storage.AddURL(shortURL, fullURL)

	var res apiRes
	baseAddr := sa.config.GetBaseAddr()
	res.Result = baseAddr + "/" + shortURL
	resBody, err := json.Marshal(&res)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
	}
	c.Set("Content-Type", "application/json")
	return c.Status(http.StatusCreated).Send(resBody)
}

func Set(app *fiber.App, sa *ServerArgs) {
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Log,
		Fields: []string{"url", "method", "latency", "status", "bytesSent"},
		Levels: []zapcore.Level{zapcore.InfoLevel},
	}))

	app.Post("/api/shorten", sa.ShortenAPI)
	app.Get("/:short", sa.ReturnURL)
	app.Post("/", sa.ShortenURL)

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusBadRequest)
	})
}

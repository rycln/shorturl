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
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap/zapcore"
)

type configer interface {
	GetBaseAddr() string
}

type storager interface {
	AddURL(string, string) bool
	GetURL(string) (string, error)
}

type fileWriter interface {
	WriteInto(*storage.StoredURL) error
}

type ServerArgs struct {
	storage    storager
	config     configer
	fileWriter fileWriter
}

func NewServerArgs(strg storager, cfg configer, fw fileWriter) *ServerArgs {
	return &ServerArgs{
		storage:    strg,
		config:     cfg,
		fileWriter: fw,
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
	ok := sa.storage.AddURL(shortURL, fullURL)
	if ok {
		surl := storage.NewStoredURL(shortURL, fullURL)
		err := sa.fileWriter.WriteInto(surl)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
	}

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
	ok := sa.storage.AddURL(shortURL, fullURL)
	if ok {
		surl := storage.NewStoredURL(shortURL, fullURL)
		err := sa.fileWriter.WriteInto(surl)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
	}

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

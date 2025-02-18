package server

import (
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/myhash"
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
	if c.Method() != http.MethodPost {
		return c.SendStatus(http.StatusBadRequest)
	}

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
	if c.Method() != http.MethodGet {
		return c.SendStatus(http.StatusBadRequest)
	}

	shortURL := c.Params("short")
	fullURL, err := sa.storage.GetURL(shortURL)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	c.Set("Location", fullURL)
	return c.SendStatus(http.StatusTemporaryRedirect)
}

func Set(app *fiber.App, sa *ServerArgs) {
	app.Use(func(c *fiber.Ctx) error {
		c.Status(http.StatusBadRequest)
		return c.Next()
	})

	app.All("/", sa.ShortenURL)
	app.All("/:short", sa.ReturnURL)
}

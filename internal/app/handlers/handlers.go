package handlers

import (
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	config "github.com/rycln/shorturl/configs"
	"github.com/rycln/shorturl/internal/app/mem"
	"github.com/rycln/shorturl/internal/app/myhash"
)

type HandlerVariables struct {
	store mem.Storager
}

func NewHandlerVariables(store mem.Storager) HandlerVariables {
	return HandlerVariables{
		store: store,
	}
}

func (hv HandlerVariables) ShortenURL(c *fiber.Ctx) error {
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
	hv.store.AddURL(shortURL, fullURL)
	c.Set("Content-Type", "text/plain")
	baseAddr := config.GetBaseAddr()
	if baseAddr == "" {
		baseAddr = "http://localhost:8080"
	}
	return c.Status(http.StatusCreated).SendString(baseAddr + "/" + shortURL)
}

func (hv HandlerVariables) ReturnURL(c *fiber.Ctx) error {
	if c.Method() != http.MethodGet {
		return c.SendStatus(http.StatusBadRequest)
	}

	shortURL := c.Params("short")
	fullURL, err := hv.store.GetURL(shortURL)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	c.Set("Location", fullURL)
	return c.SendStatus(http.StatusTemporaryRedirect)
}

package handlers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type retrieveStorager interface {
	GetOrigURL(context.Context, string) (string, error)
}

type Retrieve struct {
	strg retrieveStorager
}

func NewRetrieve(strg retrieveStorager) *Retrieve {
	return &Retrieve{
		strg: strg,
	}
}

func (r *Retrieve) RetrieveURL(c *fiber.Ctx) error {
	shortURL := c.Params("short")

	origURL, err := r.strg.GetOrigURL(c.UserContext(), shortURL)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	c.Set("Location", origURL)
	return c.SendStatus(http.StatusTemporaryRedirect)
}

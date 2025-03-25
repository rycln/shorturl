package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap"
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

func (r *Retrieve) Handle(c *fiber.Ctx) error {
	shortURL := c.Params("short")

	origURL, err := r.strg.GetOrigURL(c.UserContext(), shortURL)
	if err == nil {
		c.Set("Location", origURL)
		return c.SendStatus(http.StatusTemporaryRedirect)
	}
	if errors.Is(err, storage.ErrDeletedURL) {
		return c.SendStatus(http.StatusGone)
	}
	logger.Log.Info("path:"+c.Path()+", "+"func:GetOrigURL()",
		zap.Error(err),
	)
	return c.SendStatus(http.StatusBadRequest)
}

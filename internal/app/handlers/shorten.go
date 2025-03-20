package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap"
)

type shortenStorager interface {
	AddURL(context.Context, storage.ShortenedURL) error
	GetShortURL(context.Context, string) (string, error)
}

type shortenConfiger interface {
	GetBaseAddr() string
	GetKey() string
}

type Shorten struct {
	strg     shortenStorager
	cfg      shortenConfiger
	hashFunc func(string) string
}

func NewShorten(strg shortenStorager, cfg shortenConfiger, hashFunc func(string) string) *Shorten {
	return &Shorten{
		strg:     strg,
		cfg:      cfg,
		hashFunc: hashFunc,
	}
}

func (s *Shorten) Handle(c *fiber.Ctx) error {
	body := string(c.Body())
	_, err := url.ParseRequestURI(body)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	uid := getUserID(c, s.cfg.GetKey())
	origURL := body
	shortURL := s.hashFunc(origURL)

	baseAddr := s.cfg.GetBaseAddr()
	surl := storage.NewShortenedURL(uid, shortURL, origURL)
	err = s.strg.AddURL(c.UserContext(), surl)
	if err == nil {
		c.Set("Content-Type", "text/plain")
		return c.Status(http.StatusCreated).SendString(baseAddr + "/" + shortURL)
	}
	if errors.Is(err, storage.ErrConflict) {
		var err error
		shortURL, err = s.strg.GetShortURL(c.UserContext(), origURL)
		if err != nil {
			logger.Log.Info("path:"+c.Path()+", "+"func:GetShortURL()",
				zap.Error(err),
			)
			return c.SendStatus(http.StatusInternalServerError)
		}
		c.Set("Content-Type", "text/plain")
		return c.Status(http.StatusConflict).SendString(baseAddr + "/" + shortURL)
	}
	logger.Log.Info("path:"+c.Path()+", "+"func:AddURL()",
		zap.Error(err),
	)
	return c.SendStatus(http.StatusInternalServerError)
}

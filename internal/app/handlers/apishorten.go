package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap"
)

type apiStorager interface {
	AddURL(context.Context, storage.ShortenedURL) error
	GetShortURL(context.Context, string) (string, error)
}

type apiConfiger interface {
	GetBaseAddr() string
	GetKey() string
}

type APIShorten struct {
	strg     apiStorager
	cfg      apiConfiger
	hashFunc func(string) string
}

func NewAPIShorten(strg apiStorager, cfg apiConfiger, hashFunc func(string) string) *APIShorten {
	return &APIShorten{
		strg:     strg,
		cfg:      cfg,
		hashFunc: hashFunc,
	}
}

type apiReq struct {
	URL string `json:"url"`
}

type apiRes struct {
	Result string `json:"result"`
}

func (as *APIShorten) Handle(c *fiber.Ctx) error {
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

	uid := getUserID(c, as.cfg.GetKey())
	origURL := req.URL
	shortURL := as.hashFunc(origURL)

	var res apiRes
	baseAddr := as.cfg.GetBaseAddr()
	surl := storage.NewShortenedURL(uid, shortURL, origURL)
	err = as.strg.AddURL(c.UserContext(), surl)
	if err == nil {
		res.Result = baseAddr + "/" + shortURL
		resBody, err := json.Marshal(&res)
		if err != nil {
			logger.Log.Info("path:"+c.Path()+", "+"func:json.Marshal()",
				zap.Error(err),
			)
			c.SendStatus(http.StatusInternalServerError)
		}
		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusCreated).Send(resBody)
	}
	if errors.Is(err, storage.ErrConflict) {
		var err error
		shortURL, err = as.strg.GetShortURL(c.UserContext(), origURL)
		if err != nil {
			logger.Log.Info("path:"+c.Path()+", "+"func:GetShortURL()",
				zap.Error(err),
			)
			return c.SendStatus(http.StatusInternalServerError)
		}
		res.Result = baseAddr + "/" + shortURL
		resBody, err := json.Marshal(&res)
		if err != nil {
			logger.Log.Info("path:"+c.Path()+", "+"func:json.Marshal()",
				zap.Error(err),
			)
			c.SendStatus(http.StatusInternalServerError)
		}
		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusConflict).Send(resBody)
	}
	logger.Log.Info("path:"+c.Path()+", "+"func:AddURL()",
		zap.Error(err),
	)
	return c.SendStatus(http.StatusInternalServerError)
}

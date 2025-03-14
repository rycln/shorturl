package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap"
)

type urlAdder interface {
	AddURL(context.Context, storage.ShortenedURL) error
	AddBatchURL(context.Context, []storage.ShortenedURL) error
}

type urlGetter interface {
	GetOrigURL(context.Context, string) (string, error)
	GetShortURL(context.Context, string) (string, error)
}

type pinger interface {
	Ping(context.Context) error
}

type storager interface {
	urlAdder
	urlGetter
	pinger
}

type cfgAddresser interface {
	GetBaseAddr() string
	GetDatabaseDsn() string
}

type cfgTimeouter interface {
	TimeoutDuration() time.Duration
}

type configer interface {
	cfgAddresser
	cfgTimeouter
}

type ServerArgs struct {
	strg     storager
	cfg      configer
	hashFunc func(string) string
}

func NewServerArgs(strg storager, cfg configer, hashFunc func(string) string) *ServerArgs {
	return &ServerArgs{
		strg:     strg,
		cfg:      cfg,
		hashFunc: hashFunc,
	}
}

func (sa *ServerArgs) ShortenURL(c *fiber.Ctx) error {
	body := string(c.Body())
	_, err := url.ParseRequestURI(body)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	origURL := body
	shortURL := sa.hashFunc(origURL)

	baseAddr := sa.cfg.GetBaseAddr()
	surl := storage.NewShortenedURL(shortURL, origURL)
	err = sa.strg.AddURL(c.UserContext(), surl)
	if err == nil {
		c.Set("Content-Type", "text/plain")
		return c.Status(http.StatusCreated).SendString(baseAddr + "/" + shortURL)
	}
	if errors.Is(err, storage.ErrConflict) {
		var err error
		shortURL, err = sa.strg.GetShortURL(c.UserContext(), origURL)
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
	shortURL := sa.hashFunc(origURL)

	var res apiRes
	baseAddr := sa.cfg.GetBaseAddr()
	surl := storage.NewShortenedURL(shortURL, origURL)
	err = sa.strg.AddURL(c.UserContext(), surl)
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
		shortURL, err = sa.strg.GetShortURL(c.UserContext(), origURL)
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

type batchReq struct {
	ID      string `json:"correlation_id"`
	OrigURL string `json:"original_url"`
}

type batchRes struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

func newBatchRes(id, shortURL string) batchRes {
	return batchRes{
		ID:       id,
		ShortURL: shortURL,
	}
}

func (sa *ServerArgs) ShortenBatch(c *fiber.Ctx) error {
	if !c.Is("json") {
		return c.SendStatus(http.StatusBadRequest)
	}

	var reqBatches = []batchReq{}
	err := json.Unmarshal(c.Body(), &reqBatches)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	surls := make([]storage.ShortenedURL, len(reqBatches))
	resBatches := make([]batchRes, len(reqBatches))
	baseAddr := sa.cfg.GetBaseAddr()
	for i, b := range reqBatches {
		_, err = url.ParseRequestURI(b.OrigURL)
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}
		shortURL := sa.hashFunc(b.OrigURL)
		surls[i] = storage.NewShortenedURL(shortURL, b.OrigURL)
		resBatches[i] = newBatchRes(b.ID, baseAddr+"/"+shortURL)
	}

	err = sa.strg.AddBatchURL(c.UserContext(), surls)
	if err != nil {
		logger.Log.Info("path:"+c.Path()+", "+"func:AddBatchURL()",
			zap.Error(err),
		)
		return c.SendStatus(http.StatusInternalServerError)
	}
	resBody, err := json.Marshal(&resBatches)
	if err != nil {
		logger.Log.Info("path:"+c.Path()+", "+"func:json.Marshal()",
			zap.Error(err),
		)
		c.SendStatus(http.StatusInternalServerError)
	}
	c.Set("Content-Type", "application/json")
	return c.Status(http.StatusCreated).Send(resBody)
}

func (sa *ServerArgs) ReturnURL(c *fiber.Ctx) error {
	shortURL := c.Params("short")

	origURL, err := sa.strg.GetOrigURL(c.UserContext(), shortURL)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	c.Set("Location", origURL)
	return c.SendStatus(http.StatusTemporaryRedirect)
}

func (sa *ServerArgs) PingDB(c *fiber.Ctx) error {
	err := sa.strg.Ping(c.UserContext())
	if err != nil {
		logger.Log.Info("path:"+c.Path()+", "+"func:PingContext()",
			zap.Error(err),
		)
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}

package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/myhash"
	"github.com/rycln/shorturl/internal/app/storage"
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
		shortURL := myhash.Base62(b.OrigURL)
		surls[i] = storage.NewShortenedURL(shortURL, b.OrigURL)
		resBatches[i] = newBatchRes(b.ID, baseAddr+"/"+shortURL)
	}

	ctx, cancel := context.WithTimeout(c.Context(), 1*time.Second)
	defer cancel()
	err = sa.strg.AddURL(ctx, surls...)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	resBody, err := json.Marshal(&resBatches)
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

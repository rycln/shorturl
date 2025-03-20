package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap"
)

type batchStorager interface {
	AddBatchURL(context.Context, []storage.ShortenedURL) error
}

type batchConfiger interface {
	GetBaseAddr() string
	GetKey() string
}

type ShortenBatch struct {
	strg     batchStorager
	cfg      batchConfiger
	hashFunc func(string) string
}

func NewShortenBatch(strg batchStorager, cfg batchConfiger, hashFunc func(string) string) *ShortenBatch {
	return &ShortenBatch{
		strg:     strg,
		cfg:      cfg,
		hashFunc: hashFunc,
	}
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

func (sb *ShortenBatch) Handle(c *fiber.Ctx) error {
	if !c.Is("json") {
		return c.SendStatus(http.StatusBadRequest)
	}

	var reqBatches = []batchReq{}
	err := json.Unmarshal(c.Body(), &reqBatches)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	uid := getUserID(c, sb.cfg.GetKey())
	surls := make([]storage.ShortenedURL, len(reqBatches))
	resBatches := make([]batchRes, len(reqBatches))
	baseAddr := sb.cfg.GetBaseAddr()
	for i, b := range reqBatches {
		_, err = url.ParseRequestURI(b.OrigURL)
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}
		shortURL := sb.hashFunc(b.OrigURL)
		surls[i] = storage.NewShortenedURL(uid, shortURL, b.OrigURL)
		resBatches[i] = newBatchRes(b.ID, baseAddr+"/"+shortURL)
	}

	err = sb.strg.AddBatchURL(c.UserContext(), surls)
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

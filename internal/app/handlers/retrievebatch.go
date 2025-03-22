package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap"
)

type retrieveBatchStorager interface {
	GetAllUserURLs(context.Context, string) ([]storage.ShortenedURL, error)
}

type retrieveBatchConfiger interface {
	GetBaseAddr() string
	GetKey() string
}

type RetrieveBatch struct {
	strg retrieveBatchStorager
	cfg  retrieveBatchConfiger
}

func NewRetrieveBatch(strg retrieveBatchStorager, cfg retrieveBatchConfiger) *RetrieveBatch {
	return &RetrieveBatch{
		strg: strg,
		cfg:  cfg,
	}
}

type retBatchRes struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"original_url"`
}

func newRetBatchRes(shortURL, origURL string) retBatchRes {
	return retBatchRes{
		ShortURL: shortURL,
		OrigURL:  origURL,
	}
}

func (rb *RetrieveBatch) Handle(c *fiber.Ctx) error {
	key := rb.cfg.GetKey()
	_, uid, err := getTokenAndUID(c, key)
	if err != nil {
		uid = makeUserID()
		jwt, err := makeTokenString(uid, key)
		if err != nil {
			logger.Log.Info("path:"+c.Path()+", "+"func:makeTokenString()",
				zap.Error(err),
			)
			return c.SendStatus(http.StatusInternalServerError)
		}
		cookie := new(fiber.Cookie)
		cookie.Name = "Authorization"
		cookie.Value = fmt.Sprintf("Bearer %s", jwt)
		c.Cookie(cookie)
		return c.SendStatus(http.StatusNoContent)
	}

	surls, err := rb.strg.GetAllUserURLs(c.UserContext(), uid)
	if err == nil {
		resBatches := make([]retBatchRes, len(surls))
		baseAddr := rb.cfg.GetBaseAddr()
		for i, surl := range surls {
			resBatches[i] = newRetBatchRes(baseAddr+"/"+surl.ShortURL, surl.OrigURL)
		}
		resBody, err := json.Marshal(&resBatches)
		if err != nil {
			logger.Log.Info("path:"+c.Path()+", "+"func:json.Marshal()",
				zap.Error(err),
			)
			c.SendStatus(http.StatusInternalServerError)
		}
		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusOK).Send(resBody)
	}
	if errors.Is(err, storage.ErrNotExist) {
		return c.SendStatus(http.StatusNoContent)
	}
	logger.Log.Info("path:"+c.Path()+", "+"func:GetAllUserURLs()",
		zap.Error(err),
	)
	return c.SendStatus(http.StatusInternalServerError)
}

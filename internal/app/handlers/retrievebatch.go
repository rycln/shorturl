package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	if claims["ID"] == nil {
		logger.Log.Info("path:" + c.Path() + ", " + "user id is empty")
		return c.SendStatus(http.StatusInternalServerError)
	}
	uid := claims["ID"].(string)

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
		return c.Status(http.StatusCreated).Send(resBody)
	}
	if errors.Is(err, storage.ErrNoUserURLs) {
		return c.SendStatus(http.StatusNoContent)
	}
	logger.Log.Info("path:"+c.Path()+", "+"func:GetAllUserURLs()",
		zap.Error(err),
	)
	return c.SendStatus(http.StatusInternalServerError)
}

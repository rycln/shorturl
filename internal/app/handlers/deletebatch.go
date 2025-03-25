package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/logger"
	"github.com/rycln/shorturl/internal/app/storage"
	"go.uber.org/zap"
)

const batchDeletingPeriod = 10 * time.Second

type deleteBatchStorager interface {
	DeleteUserURLs(context.Context, []storage.DelShortURLs) error
}

type deleteBatchConfiger interface {
	GetBaseAddr() string
	GetKey() string
}

type DeleteBatch struct {
	strg    deleteBatchStorager
	cfg     deleteBatchConfiger
	delChan chan storage.DelShortURLs
}

func NewDeleteBatch(strg deleteBatchStorager, cfg deleteBatchConfiger) *DeleteBatch {
	delb := &DeleteBatch{
		strg:    strg,
		cfg:     cfg,
		delChan: make(chan storage.DelShortURLs, 1024),
	}
	go delb.deleteBatchProcessing()
	return delb
}

func (delb *DeleteBatch) Handle(c *fiber.Ctx) error {
	key := delb.cfg.GetKey()
	_, uid, err := getTokenAndUID(c, key)
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	var shortURLs []string
	err = json.Unmarshal(c.Body(), &shortURLs)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	delb.batchGen(uid, shortURLs)

	return c.SendStatus(http.StatusAccepted)
}

func (delb *DeleteBatch) batchGen(uid string, shortURLs []string) {
	go func() {
		for _, shortURL := range shortURLs {
			delb.delChan <- storage.NewDelShortURLs(uid, shortURL)
		}
	}()
}

func (delb *DeleteBatch) deleteBatchProcessing() {
	tick := time.NewTicker(batchDeletingPeriod)

	var shortURLBatch []storage.DelShortURLs

	for {
		select {
		case shortURL := <-delb.delChan:
			shortURLBatch = append(shortURLBatch, shortURL)
		case <-tick.C:
			if len(shortURLBatch) == 0 {
				continue
			}
			//тут отправка пачки в бд
			err := delb.strg.DeleteUserURLs(context.TODO(), shortURLBatch)
			if err != nil {
				logger.Log.Info("Cannot delete batch",
					zap.Error(err),
				)
				continue
			}

			shortURLBatch = nil
		}
	}
}

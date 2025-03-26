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

const batchDeletingPeriod = 5 * time.Second

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
	inChans chan chan storage.DelShortURLs
}

func NewDeleteBatch(strg deleteBatchStorager, cfg deleteBatchConfiger) *DeleteBatch {
	delb := &DeleteBatch{
		strg:    strg,
		cfg:     cfg,
		inChans: make(chan chan storage.DelShortURLs),
	}

	deleteBatchInit(delb)

	return delb
}

func deleteBatchInit(delb *DeleteBatch) {
	outChan := mergeChans(delb.inChans)
	go delb.deleteBatchWorker(outChan)
}

func (delb *DeleteBatch) Handle(c *fiber.Ctx) error {
	key := delb.cfg.GetKey()
	_, uid, err := getTokenAndUID(c, key)
	if err != nil {
		uid = makeUserID()
	}

	var shortURLs []string
	err = json.Unmarshal(c.Body(), &shortURLs)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	delb.makeChan(uid, shortURLs)

	return c.SendStatus(http.StatusAccepted)
}

func (delb *DeleteBatch) makeChan(uid string, shortURLs []string) {
	ch := make(chan storage.DelShortURLs)
	go func() {
		defer close(ch)
		for _, shortURL := range shortURLs {
			ch <- storage.NewDelShortURLs(uid, shortURL)
		}
	}()

	delb.inChans <- ch
}

func mergeChans(inputs <-chan chan storage.DelShortURLs) chan storage.DelShortURLs {
	out := make(chan storage.DelShortURLs)

	go func() {
		for ch := range inputs {
			go func(ch <-chan storage.DelShortURLs) {
				for dShortURL := range ch {
					out <- dShortURL
				}
			}(ch)
		}
	}()

	return out
}

func (delb *DeleteBatch) deleteBatchWorker(input <-chan storage.DelShortURLs) {
	tick := time.NewTicker(batchDeletingPeriod)

	var delShortURLBatch []storage.DelShortURLs

	for {
		select {
		case delShortURL := <-input:
			delShortURLBatch = append(delShortURLBatch, delShortURL)
		case <-tick.C:
			if len(delShortURLBatch) == 0 {
				continue
			}
			//тут отправка пачки в бд
			err := delb.strg.DeleteUserURLs(context.Background(), delShortURLBatch)
			if err != nil {
				logger.Log.Info("Cannot delete batch",
					zap.Error(err),
				)
				continue
			}

			delShortURLBatch = nil
		}
	}
}

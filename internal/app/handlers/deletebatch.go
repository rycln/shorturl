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

const (
	batchDeletingPeriod  = 10 * time.Second
	batchDeletingTimeout = 3 * time.Second
)

type deleteBatchStorager interface {
	DeleteUserURLs(context.Context, []storage.DelShortURLs) error
}

type deleteBatchConfiger interface {
	GetBaseAddr() string
	GetKey() string
}

type DeleteBatch struct {
	strg  deleteBatchStorager
	cfg   deleteBatchConfiger
	chans chan chan storage.DelShortURLs
}

func NewDeleteBatchHandler(strg deleteBatchStorager, cfg deleteBatchConfiger) func(*fiber.Ctx) error {
	delb := &DeleteBatch{
		strg:  strg,
		cfg:   cfg,
		chans: make(chan chan storage.DelShortURLs),
	}

	deleteBatchInit(delb)

	return delb.handle
}

func deleteBatchInit(delb *DeleteBatch) {
	mergedChan := mergeChans(delb.chans)
	go delb.deleteBatchWorker(mergedChan)
}

func mergeChans(inChans <-chan chan storage.DelShortURLs) chan storage.DelShortURLs {
	outChan := make(chan storage.DelShortURLs)

	go func() {
		for ch := range inChans {
			go func(ch <-chan storage.DelShortURLs) {
				for dShortURL := range ch {
					outChan <- dShortURL
				}
			}(ch)
		}
	}()

	return outChan
}

func (delb *DeleteBatch) deleteBatchWorker(inChan <-chan storage.DelShortURLs) {
	tick := time.NewTicker(batchDeletingPeriod)

	var delShortURLBatch []storage.DelShortURLs

	for {
		select {
		case delShortURL := <-inChan:
			delShortURLBatch = append(delShortURLBatch, delShortURL)
		case <-tick.C:
			if len(delShortURLBatch) == 0 {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), batchDeletingTimeout)
			err := delb.strg.DeleteUserURLs(ctx, delShortURLBatch)
			if err != nil {
				logger.Log.Info("Cannot delete batch",
					zap.Error(err),
				)
				continue
			}
			cancel()

			delShortURLBatch = nil
		}
	}
}

func (delb *DeleteBatch) handle(c *fiber.Ctx) error {
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

	delb.chans <- ch
}

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
	GetTimeoutDuration() time.Duration
	GetKey() string
}

type DeleteBatch struct {
	doneCtx context.Context
	strg    deleteBatchStorager
	cfg     deleteBatchConfiger
	chans   chan chan storage.DelShortURLs
}

func NewDeleteBatchHandler(ctx context.Context, strg deleteBatchStorager, cfg deleteBatchConfiger) func(*fiber.Ctx) error {
	delb := &DeleteBatch{
		doneCtx: ctx,
		strg:    strg,
		cfg:     cfg,
		chans:   make(chan chan storage.DelShortURLs),
	}

	deleteBatchInit(delb)

	return delb.handle
}

func deleteBatchInit(delb *DeleteBatch) {
	mergedChan := mergeChans(delb.doneCtx, delb.chans)
	go delb.deleteBatchWorker(mergedChan)
}

func mergeChans(ctx context.Context, inChans <-chan chan storage.DelShortURLs) chan storage.DelShortURLs {
	outChan := make(chan storage.DelShortURLs)

	go func() {
		for ch := range inChans {
			select {
			case <-ctx.Done():
				return
			default:
				go func(ch <-chan storage.DelShortURLs) {
					for dShortURL := range ch {
						select {
						case <-ctx.Done():
							return
						default:
							outChan <- dShortURL
						}
					}
				}(ch)
			}
		}
	}()

	return outChan
}

func (delb *DeleteBatch) deleteBatchWorker(inChan <-chan storage.DelShortURLs) {
	tick := time.NewTicker(batchDeletingPeriod)

	var delShortURLBatch []storage.DelShortURLs

	for {
		select {
		case <-delb.doneCtx.Done():
			return
		case delShortURL := <-inChan:
			delShortURLBatch = append(delShortURLBatch, delShortURL)
		case <-tick.C:
			if len(delShortURLBatch) == 0 {
				continue
			}

			ctx, cancel := context.WithTimeout(delb.doneCtx, delb.cfg.GetTimeoutDuration())
			err := delb.strg.DeleteUserURLs(ctx, delShortURLBatch)
			cancel()
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
			select {
			case <-delb.doneCtx.Done():
				return
			default:
				ch <- storage.NewDelShortURLs(uid, shortURL)
			}
		}
	}()

	delb.chans <- ch
}

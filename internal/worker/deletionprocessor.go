// Package worker implements a background URL deletion processor for the URL shortener service.
package worker

import (
	"context"
	"time"

	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type batchDeleteServicer interface {
	DeleteURLsBatch(context.Context, []*models.DelURLReq) error
}

// DeletionProcessor is a background worker that processes URL deletions in batches.
//
// The processor collects deletion requests in memory and flushes them to storage
// either when batch size limit is reached or on timer expiration.
type DeletionProcessor struct {
	ctx                context.Context
	cancel             context.CancelFunc
	batchDeleteService batchDeleteServicer
	delChans           chan chan *models.DelURLReq
}

// NewDeletionProcessor creates new processor instance.
func NewDeletionProcessor(batchDeleteService batchDeleteServicer) *DeletionProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	return &DeletionProcessor{
		ctx:                ctx,
		cancel:             cancel,
		batchDeleteService: batchDeleteService,
		delChans:           make(chan chan *models.DelURLReq, 10),
	}
}

// Shutdown stops the processor.
func (p *DeletionProcessor) Shutdown() {
	p.cancel()
	close(p.delChans)
}

// Run starts the background processing loop.
//
// The processor will handle requests until Shutdown() is called.
func (p *DeletionProcessor) Run(period time.Duration, timeout time.Duration) chan struct{} {
	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		inChan := fanIn(p.ctx, p.delChans)

		tick := time.NewTicker(period)

		var delBatch []*models.DelURLReq

		for {
			select {
			case <-p.ctx.Done():
				return
			case durl := <-inChan:
				delBatch = append(delBatch, durl)
			case <-tick.C:
				if len(delBatch) == 0 {
					continue
				}

				ctx, cancel := context.WithTimeout(p.ctx, timeout)
				err := p.batchDeleteService.DeleteURLsBatch(ctx, delBatch)
				if err != nil {
					logger.Log.Info("Cannot delete batch", zap.Error(err))
					cancel()
					continue
				}
				cancel()
				delBatch = nil
			}
		}
	}()

	return doneCh
}

// AddURLsIntoDeletionQueue enqueues URLs for deletion.
//
// Non-blocking method for use in HTTP handlers.
func (p *DeletionProcessor) AddURLsIntoDeletionQueue(uid models.UserID, shorts []models.ShortURL) {
	delCh := make(chan *models.DelURLReq)

	go func() {
		defer close(delCh)
		for _, short := range shorts {
			select {
			case <-p.ctx.Done():
				return
			default:
				delCh <- &models.DelURLReq{
					UID:   uid,
					Short: short,
				}
			}
		}
	}()

	p.delChans <- delCh
}

func fanIn(ctx context.Context, inChans <-chan chan *models.DelURLReq) chan *models.DelURLReq {
	outChan := make(chan *models.DelURLReq)

	go func() {
		defer close(outChan)

		for ch := range inChans {
			select {
			case <-ctx.Done():
				return
			default:
				go func(ch <-chan *models.DelURLReq) {
					for durl := range ch {
						select {
						case <-ctx.Done():
							return
						default:
							outChan <- durl
						}
					}
				}(ch)
			}
		}
	}()

	return outChan
}

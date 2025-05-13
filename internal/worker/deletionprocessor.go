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

type DeletionProcessor struct {
	ctx                context.Context
	cancel             context.CancelFunc
	batchDeleteService batchDeleteServicer
	delChans           chan chan *models.DelURLReq
}

func NewDeletionProcessor(batchDeleteService batchDeleteServicer) *DeletionProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	return &DeletionProcessor{
		ctx:                ctx,
		cancel:             cancel,
		batchDeleteService: batchDeleteService,
		delChans:           make(chan chan *models.DelURLReq, 10),
	}
}

func (p *DeletionProcessor) Shutdown() {
	p.cancel()
	close(p.delChans)
}

func (p *DeletionProcessor) Run(period time.Duration, timeout time.Duration) {
	go func() {
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
}

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

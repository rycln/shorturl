package service

import (
	"context"

	"github.com/rycln/shorturl/internal/models"
)

type BatchDeleteStorage interface {
	DeleteRequestedURLs(context.Context, []models.DelURLReq) error
}

type BatchDelete struct {
	strg     BatchDeleteStorage
	delChans chan chan *models.DelURLReq
	cancelCh chan struct{}
}

func NewBatchDelete(strg BatchDeleteStorage, delChans chan chan *models.DelURLReq, cancelCh chan struct{}) *BatchDelete {
	return &BatchDelete{
		strg:     strg,
		delChans: delChans,
		cancelCh: cancelCh,
	}
}

func (s *BatchDelete) UserURLsAsyncDeletion(uid models.UserID, shorts []models.ShortURL) {
	delCh := make(chan *models.DelURLReq)

	go func() {
		defer close(delCh)
		for _, short := range shorts {
			select {
			case <-s.cancelCh:
				return
			default:
				delCh <- &models.DelURLReq{
					UID:   uid,
					Short: short,
				}
			}
		}
	}()

	s.delChans <- delCh
}

func (s *BatchDelete) DeleteURLsBatch(ctx context.Context, urls []models.DelURLReq) error {
	err := s.strg.DeleteRequestedURLs(ctx, urls)
	if err != nil {
		return err
	}
	return nil
}

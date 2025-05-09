package service

import (
	"context"

	"github.com/rycln/shorturl/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type BatchDeleterStorage interface {
	DeleteRequestedURLs(context.Context, []models.DelURLReq) error
}

type BatchDeleter struct {
	strg     BatchDeleterStorage
	delChans chan chan *models.DelURLReq
	cancelCh chan struct{}
}

func NewBatchDeleter(strg BatchDeleterStorage, delChans chan chan *models.DelURLReq, cancelCh chan struct{}) *BatchDeleter {
	return &BatchDeleter{
		strg:     strg,
		delChans: delChans,
		cancelCh: cancelCh,
	}
}

func (s *BatchDeleter) UserURLsAsyncDeletion(uid models.UserID, shorts []models.ShortURL) {
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

func (s *BatchDeleter) DeleteURLsBatch(ctx context.Context, urls []models.DelURLReq) error {
	err := s.strg.DeleteRequestedURLs(ctx, urls)
	if err != nil {
		return err
	}
	return nil
}

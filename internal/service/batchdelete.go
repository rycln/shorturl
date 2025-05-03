package service

import "github.com/rycln/shorturl/internal/models"

type BatchDelete struct {
	delChans chan chan *models.DelURLReq
	cancelCh chan struct{}
}

func NewBatchDelete(delChans chan chan *models.DelURLReq, cancelCh chan struct{}) *BatchDelete {
	return &BatchDelete{
		delChans: delChans,
		cancelCh: cancelCh,
	}
}

func (s *BatchDelete) DeleteUserURLsByShorts(uid models.UserID, shorts []models.ShortURL) {
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

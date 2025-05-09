package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/models"
	"github.com/rycln/shorturl/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBatchDelete_UserURLsAsyncDeletion(t *testing.T) {

	t.Run("valid test", func(t *testing.T) {
		s := &BatchDeleter{
			delChans: make(chan chan *models.DelURLReq, 1),
			cancelCh: make(chan struct{}),
		}

		urls := []models.ShortURL{"url1", "url2", "url3"}

		s.UserURLsAsyncDeletion(testUserID, urls)

		delCh := <-s.delChans

		var received []*models.DelURLReq
		for req := range delCh {
			received = append(received, req)
		}

		assert.Equal(t, len(urls), len(received))

		for i, req := range received {
			assert.Equal(t, testUserID, req.UID)
			assert.Equal(t, urls[i], req.Short)
		}
	})

	t.Run("cancellation", func(t *testing.T) {
		s := &BatchDeleter{
			delChans: make(chan chan *models.DelURLReq, 1),
			cancelCh: make(chan struct{}),
		}

		urls := make([]models.ShortURL, 1000)
		for i := range urls {
			urls[i] = models.ShortURL(fmt.Sprintf("url%d", i))
		}

		go func() {
			time.Sleep(10 * time.Microsecond)
			close(s.cancelCh)
		}()

		s.UserURLsAsyncDeletion(testUserID, urls)
		delCh := <-s.delChans

		count := 0
		for range delCh {
			count++
		}

		assert.NotEqual(t, 0, count)
		assert.NotEqual(t, len(urls), count, "cancel did not work")
	})
}

func TestBatchDelete_DeleteURLsBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	delChans := make(chan chan *models.DelURLReq, 1)
	cancelCh := make(chan struct{})
	mStrg := mocks.NewMockBatchDeleterStorage(ctrl)

	s := NewBatchDeleter(mStrg, delChans, cancelCh)

	testDelReqs := []models.DelURLReq{
		testDelReq,
	}

	t.Run("valid test", func(t *testing.T) {
		mStrg.EXPECT().DeleteRequestedURLs(context.Background(), testDelReqs).Return(nil)

		err := s.DeleteURLsBatch(context.Background(), testDelReqs)
		assert.NoError(t, err)
	})

	t.Run("some error", func(t *testing.T) {
		mStrg.EXPECT().DeleteRequestedURLs(context.Background(), testDelReqs).Return(errTest)

		err := s.DeleteURLsBatch(context.Background(), testDelReqs)
		assert.Error(t, err)
	})
}

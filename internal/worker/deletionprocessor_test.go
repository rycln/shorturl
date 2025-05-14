package worker

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/models"
	"github.com/rycln/shorturl/internal/worker/mocks"
)

const (
	testUserID  = models.UserID("1")
	testTicker  = time.Duration(10) * time.Millisecond
	testTimeout = time.Duration(10) * time.Millisecond
)

var errTest = errors.New("test error")

func TestDeletionProcessor_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("valid test", func(t *testing.T) {
		mServ := mocks.NewMockbatchDeleteServicer(ctrl)

		var wg sync.WaitGroup
		wg.Add(1)

		mServ.EXPECT().DeleteURLsBatch(gomock.Any(), gomock.Any()).Return(nil).Times(1).Do(func(_, _ interface{}) {
			wg.Done()
		})

		p := NewDeletionProcessor(mServ)
		defer p.Shutdown()

		urls := []models.ShortURL{"url1", "url2", "url3"}

		p.Run(testTicker, testTimeout)
		p.AddURLsIntoDeletionQueue(testUserID, urls)

		var done = make(chan struct{})

		go func() {
			wg.Wait()
			close(done)
		}()

		<-done
	})

	t.Run("serv error", func(t *testing.T) {
		mServ := mocks.NewMockbatchDeleteServicer(ctrl)

		var wg sync.WaitGroup
		wg.Add(1)

		mServ.EXPECT().DeleteURLsBatch(gomock.Any(), gomock.Any()).Return(errTest).Times(1).Do(func(_, _ interface{}) {
			wg.Done()
		})

		p := NewDeletionProcessor(mServ)
		defer p.Shutdown()

		urls := []models.ShortURL{"url1", "url2", "url3"}

		p.Run(testTicker, testTimeout)
		p.AddURLsIntoDeletionQueue(testUserID, urls)

		var done = make(chan struct{})

		go func() {
			wg.Wait()
			close(done)
		}()

		<-done
	})

	t.Run("shutdown", func(t *testing.T) {
		mServ := mocks.NewMockbatchDeleteServicer(ctrl)

		p := NewDeletionProcessor(mServ)

		urls := []models.ShortURL{"url1", "url2", "url3"}

		p.Run(testTicker, testTimeout)
		p.AddURLsIntoDeletionQueue(testUserID, urls)

		p.Shutdown()
	})
}

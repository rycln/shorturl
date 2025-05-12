package worker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rycln/shorturl/internal/models"
	"github.com/stretchr/testify/assert"
)

const (
	testUserID = models.UserID("1")
)

func TestDeletionProcessor_AddURLsIntoDeletionQueue(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := &DeletionProcessor{
			ctx:      ctx,
			cancel:   cancel,
			delChans: make(chan chan *models.DelURLReq, 10),
		}

		urls := []models.ShortURL{"url1", "url2", "url3"}

		p.AddURLsIntoDeletionQueue(testUserID, urls)

		delCh := <-p.delChans

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
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := &DeletionProcessor{
			ctx:      ctx,
			cancel:   cancel,
			delChans: make(chan chan *models.DelURLReq, 10),
		}

		urls := make([]models.ShortURL, 1000)
		for i := range urls {
			urls[i] = models.ShortURL(fmt.Sprintf("url%d", i))
		}

		go func() {
			time.Sleep(10 * time.Microsecond)
			p.cancel()
		}()

		p.AddURLsIntoDeletionQueue(testUserID, urls)
		delCh := <-p.delChans

		count := 0
		for range delCh {
			count++
		}

		assert.NotEqual(t, 0, count)
		assert.NotEqual(t, len(urls), count, "cancel did not work")
	})
}

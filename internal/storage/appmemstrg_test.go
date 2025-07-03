package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/rycln/shorturl/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppMemStorage_AddURLPair(t *testing.T) {
	strg := NewAppMemStorage()

	umap := make(map[models.ShortURL]models.OrigURL)
	umap[testShortURL] = testOrigURL
	strg.pairs[testUserID] = umap

	t.Run("valid test", func(t *testing.T) {
		pair := models.URLPair{
			UID:   "123",
			Short: "234",
			Orig:  "https://ya.ru/123",
		}
		err := strg.AddURLPair(context.Background(), &pair)
		assert.NoError(t, err)
	})

	t.Run("valid test #2", func(t *testing.T) {
		pair := models.URLPair{
			UID:   testUserID,
			Short: "345",
			Orig:  "https://ya.ru/123",
		}
		err := strg.AddURLPair(context.Background(), &pair)
		assert.NoError(t, err)
	})

	t.Run("ctx expired", func(t *testing.T) {
		pair := models.URLPair{
			UID:   "123",
			Short: "234",
			Orig:  "https://ya.ru/123",
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := strg.AddURLPair(ctx, &pair)
		assert.Error(t, err)
	})

	t.Run("conflict", func(t *testing.T) {
		pair := models.URLPair{
			UID:   testUserID,
			Short: testShortURL,
			Orig:  testOrigURL,
		}
		err := strg.AddURLPair(context.Background(), &pair)
		assert.ErrorIs(t, err, errConflict)
	})

	t.Run("conflict #2", func(t *testing.T) {
		pair := models.URLPair{
			UID:   testOtherUserID,
			Short: testShortURL,
			Orig:  testOrigURL,
		}
		err := strg.AddURLPair(context.Background(), &pair)
		assert.ErrorIs(t, err, errConflict)
	})
}

func BenchmarkAppMemStorage_AddURLPair(b *testing.B) {
	b.Run("add unique pair", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			storage := NewAppMemStorage()

			err := storage.AddURLPair(context.Background(), &testPair)
			require.NoError(b, err)
		}
	})

	b.Run("add 100 unique pairs", func(b *testing.B) {
		len := 100
		pairs := make([]models.URLPair, len)
		for i := range len {
			pairs[i] = models.URLPair{
				UID:   models.UserID(fmt.Sprintf("user-%d", i)),
				Orig:  models.OrigURL(fmt.Sprintf("https://site.com/page%d", i)),
				Short: models.ShortURL(fmt.Sprintf("hash-%d", i)),
			}
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage := NewAppMemStorage()

			for j := range len {
				err := storage.AddURLPair(context.Background(), &pairs[j])
				require.NoError(b, err)
			}
		}
	})

	b.Run("add 1000 unique pairs", func(b *testing.B) {
		len := 1000
		pairs := make([]models.URLPair, len)
		for i := range len {
			pairs[i] = models.URLPair{
				UID:   models.UserID(fmt.Sprintf("user-%d", i)),
				Orig:  models.OrigURL(fmt.Sprintf("https://site.com/page%d", i)),
				Short: models.ShortURL(fmt.Sprintf("hash-%d", i)),
			}
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage := NewAppMemStorage()

			for j := range len {
				err := storage.AddURLPair(context.Background(), &pairs[j])
				require.NoError(b, err)
			}
		}
	})

	b.Run("add duplicate pair", func(b *testing.B) {
		storage := NewAppMemStorage()
		err := storage.AddURLPair(context.Background(), &testPair)
		require.NoError(b, err)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := storage.AddURLPair(context.Background(), &testPair)
			require.Error(b, err)
		}
	})
}

func TestAppMemStorage_GetURLPairByShort(t *testing.T) {
	strg := NewAppMemStorage()

	umap := make(map[models.ShortURL]models.OrigURL)
	umap[testShortURL] = testOrigURL
	strg.pairs[testUserID] = umap
	strg.deleted[testDeletedShort] = struct{}{}

	t.Run("valid test", func(t *testing.T) {
		pair, err := strg.GetURLPairByShort(context.Background(), testShortURL)
		assert.NoError(t, err)
		assert.Equal(t, testPair, *pair)
	})

	t.Run("ctx expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := strg.GetURLPairByShort(ctx, testShortURL)
		assert.Error(t, err)
	})

	t.Run("deleted url error", func(t *testing.T) {
		_, err := strg.GetURLPairByShort(context.Background(), testDeletedShort)
		assert.ErrorIs(t, err, errDeletedURL)
	})

	t.Run("not exist error", func(t *testing.T) {
		_, err := strg.GetURLPairByShort(context.Background(), models.ShortURL("not exist"))
		assert.ErrorIs(t, err, errNotExist)
	})
}

func BenchmarkAppMemStorage_GetURLPairByShort(b *testing.B) {
	b.Run("get pair", func(b *testing.B) {
		storage := NewAppMemStorage()
		_ = storage.AddURLPair(context.Background(), &testPair)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := storage.GetURLPairByShort(context.Background(), testShortURL)
			require.NoError(b, err)
		}
	})
}

func TestAppMemStorage_AddBatchURLPairs(t *testing.T) {
	strg := NewAppMemStorage()

	pairs := []models.URLPair{
		{
			UID:   testUserID,
			Short: testShortURL,
			Orig:  testOrigURL,
		},
		{
			UID:   testUserID,
			Short: "132",
			Orig:  "https://ya.ru/123",
		},
	}

	t.Run("valid test", func(t *testing.T) {
		err := strg.AddBatchURLPairs(context.Background(), pairs)
		assert.NoError(t, err)
	})

	t.Run("ctx expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := strg.AddBatchURLPairs(ctx, pairs)
		assert.Error(t, err)
	})
}

func BenchmarkAppMemStorage_AddBatchURLPairs(b *testing.B) {
	b.Run("add pair", func(b *testing.B) {
		storage := NewAppMemStorage()
		pair := make([]models.URLPair, 1)
		pair[0] = testPair
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := storage.AddBatchURLPairs(context.Background(), pair)
			require.NoError(b, err)
		}
	})

	b.Run("add 100 pairs", func(b *testing.B) {
		storage := NewAppMemStorage()
		len := 100
		pairs := make([]models.URLPair, len)
		for i := range len {
			pairs[i] = models.URLPair{
				UID:   models.UserID(fmt.Sprintf("user-%d", i)),
				Orig:  models.OrigURL(fmt.Sprintf("https://site.com/page%d", i)),
				Short: models.ShortURL(fmt.Sprintf("hash-%d", i)),
			}
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := storage.AddBatchURLPairs(context.Background(), pairs)
			require.NoError(b, err)
		}
	})

	b.Run("add 1000 pairs", func(b *testing.B) {
		storage := NewAppMemStorage()
		len := 1000
		pairs := make([]models.URLPair, len)
		for i := range len {
			pairs[i] = models.URLPair{
				UID:   models.UserID(fmt.Sprintf("user-%d", i)),
				Orig:  models.OrigURL(fmt.Sprintf("https://site.com/page%d", i)),
				Short: models.ShortURL(fmt.Sprintf("hash-%d", i)),
			}
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := storage.AddBatchURLPairs(context.Background(), pairs)
			require.NoError(b, err)
		}
	})
}

func TestAppMemStorage_GetURLPairBatchByUserID(t *testing.T) {
	strg := NewAppMemStorage()

	t.Run("valid test", func(t *testing.T) {
		umap := make(map[models.ShortURL]models.OrigURL)
		umap[testShortURL] = testOrigURL
		strg.pairs[testUserID] = umap

		pairs, err := strg.GetURLPairBatchByUserID(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, testPair, pairs[0])
	})

	t.Run("ctx expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := strg.GetURLPairBatchByUserID(ctx, testUserID)
		assert.Error(t, err)
	})

	t.Run("not exist error", func(t *testing.T) {
		_, err := strg.GetURLPairBatchByUserID(context.Background(), "user id")
		assert.ErrorIs(t, err, errNotExist)
	})
}

func BenchmarkAppMemStorage_GetURLPairBatchByUserID(b *testing.B) {
	b.Run("get url", func(b *testing.B) {
		storage := NewAppMemStorage()
		userPairs := make(map[models.ShortURL]models.OrigURL)
		userPairs[testShortURL] = testOrigURL
		storage.pairs[testUserID] = userPairs
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := storage.GetURLPairBatchByUserID(context.Background(), testUserID)
			require.NoError(b, err)
		}
	})

	b.Run("get 100 urls", func(b *testing.B) {
		storage := NewAppMemStorage()
		userPairs := make(map[models.ShortURL]models.OrigURL)
		for i := 0; i < 100; i++ {
			userPairs[models.ShortURL(fmt.Sprintf("hash-%d", i))] = models.OrigURL(fmt.Sprintf("https://site.com/page%d", i))
		}
		storage.pairs[testUserID] = userPairs
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := storage.GetURLPairBatchByUserID(context.Background(), testUserID)
			require.NoError(b, err)
		}
	})

	b.Run("get 1000 urls", func(b *testing.B) {
		storage := NewAppMemStorage()
		userPairs := make(map[models.ShortURL]models.OrigURL)
		for i := 0; i < 1000; i++ {
			userPairs[models.ShortURL(fmt.Sprintf("hash-%d", i))] = models.OrigURL(fmt.Sprintf("https://site.com/page%d", i))
		}
		storage.pairs[testUserID] = userPairs
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := storage.GetURLPairBatchByUserID(context.Background(), testUserID)
			require.NoError(b, err)
		}
	})
}

func TestAppMemStorage_DeleteRequestedURLs(t *testing.T) {
	strg := NewAppMemStorage()

	delurls := []*models.DelURLReq{
		{
			UID:   testUserID,
			Short: testShortURL,
		},
	}

	t.Run("valid test", func(t *testing.T) {
		err := strg.DeleteRequestedURLs(context.Background(), delurls)
		assert.NoError(t, err)
	})

	t.Run("ctx expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := strg.DeleteRequestedURLs(ctx, delurls)
		assert.Error(t, err)
	})
}

func BenchmarkAppMemStorage_DeleteRequestedURLs(b *testing.B) {
	b.Run("delete request", func(b *testing.B) {
		storage := NewAppMemStorage()

		delURLs := []*models.DelURLReq{
			&testDelReq,
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := storage.DeleteRequestedURLs(context.Background(), delURLs)
			require.NoError(b, err)
		}
	})

	b.Run("100 delete requests", func(b *testing.B) {
		storage := NewAppMemStorage()

		len := 100
		delURLs := make([]*models.DelURLReq, len)
		for i := range len {
			delURLs[i] = &models.DelURLReq{
				UID:   models.UserID(fmt.Sprintf("user-%d", i)),
				Short: models.ShortURL(fmt.Sprintf("hash-%d", i)),
			}
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := storage.DeleteRequestedURLs(context.Background(), delURLs)
			require.NoError(b, err)
		}
	})

	b.Run("1000 delete requests", func(b *testing.B) {
		storage := NewAppMemStorage()

		len := 1000
		delURLs := make([]*models.DelURLReq, len)
		for i := range len {
			delURLs[i] = &models.DelURLReq{
				UID:   models.UserID(fmt.Sprintf("user-%d", i)),
				Short: models.ShortURL(fmt.Sprintf("hash-%d", i)),
			}
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := storage.DeleteRequestedURLs(context.Background(), delURLs)
			require.NoError(b, err)
		}
	})
}

func TestAppMemStorage_Ping(t *testing.T) {
	strg := NewAppMemStorage()

	t.Run("valid test", func(t *testing.T) {
		err := strg.Ping(context.Background())
		assert.NoError(t, err)
	})
}

func BenchmarkAppMemStorage_Ping(b *testing.B) {
	b.Run("ping", func(b *testing.B) {
		storage := NewAppMemStorage()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := storage.Ping(context.Background())
			require.NoError(b, err)
		}
	})
}

func TestAppMemStorage_GetStat(t *testing.T) {
	strg := NewAppMemStorage()

	umap := make(map[models.ShortURL]models.OrigURL)
	umap[testShortURL] = testOrigURL
	strg.pairs[testUserID] = umap

	users := 1
	urls := 1

	t.Run("valid test", func(t *testing.T) {
		stat, err := strg.GetStats(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, urls, stat.URLs)
		assert.Equal(t, users, stat.Users)
	})

	t.Run("canceled ctx", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := strg.GetStats(ctx)
		assert.Error(t, err)
	})
}

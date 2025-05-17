package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/rycln/shorturl/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_AddURLPair(t *testing.T) {
	strg, err := NewFileStorage(testFileName)
	require.NoError(t, err)
	defer os.Remove(strg.strgFileName)
	defer os.Remove(strg.delFileName)

	t.Run("valid test", func(t *testing.T) {
		err = strg.AddURLPair(context.Background(), &testPair)
		assert.NoError(t, err)
	})

	t.Run("ctx expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = strg.AddURLPair(ctx, &testPair)
		assert.Error(t, err)
	})

	t.Run("conflict error", func(t *testing.T) {
		pair := &models.URLPair{
			UID:   testUserID,
			Short: "1234",
			Orig:  "https://ya.ru/123",
		}

		file, err := os.OpenFile(strg.strgFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		require.NoError(t, err)
		defer file.Close()

		enc := json.NewEncoder(file)
		require.NoError(t, err)
		err = enc.Encode(pair)
		require.NoError(t, err)

		err = strg.AddURLPair(context.Background(), pair)
		assert.ErrorIs(t, err, ErrConflict)
	})
}

func BenchmarkFileStorage_AddURLPair(b *testing.B) {
	b.Run("add unique pair", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage.AddURLPair(context.Background(), &testPair)
		}
	})

	b.Run("add 100 unique pairs", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

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
			for j := range len {
				storage.AddURLPair(context.Background(), &pairs[j])
			}
		}
	})

	b.Run("add 1000 unique pairs", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

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
			for j := range len {
				storage.AddURLPair(context.Background(), &pairs[j])
			}
		}
	})

	b.Run("add duplicate pair", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

		storage.AddURLPair(context.Background(), &testPair)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage.AddURLPair(context.Background(), &testPair)
		}
	})
}

func TestFileStorage_GetURLPairByShort(t *testing.T) {
	strg, err := NewFileStorage(testFileName)
	require.NoError(t, err)
	defer os.Remove(strg.strgFileName)
	defer os.Remove(strg.delFileName)

	file, err := os.OpenFile(strg.strgFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	require.NoError(t, err)
	defer file.Close()

	enc := json.NewEncoder(file)
	require.NoError(t, err)
	err = enc.Encode(&testPair)
	require.NoError(t, err)

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
		file, err := os.OpenFile(strg.delFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		require.NoError(t, err)
		defer file.Close()

		enc := json.NewEncoder(file)
		require.NoError(t, err)
		err = enc.Encode(&testDelReq)
		require.NoError(t, err)

		_, err = strg.GetURLPairByShort(context.Background(), testDeletedShort)
		assert.ErrorIs(t, err, ErrDeletedURL)
	})
}

func BenchmarkFileStorage_GetURLPairByShort(b *testing.B) {
	b.Run("get pair", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

		_ = storage.AddURLPair(context.Background(), &testPair)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage.GetURLPairByShort(context.Background(), testShortURL)
		}
	})
}

func TestFileStorage_AddBatchURLPairs(t *testing.T) {
	strg, err := NewFileStorage(testFileName)
	require.NoError(t, err)
	defer os.Remove(strg.strgFileName)
	defer os.Remove(strg.delFileName)

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

func BenchmarkFileStorage_AddBatchURLPairs(b *testing.B) {
	b.Run("add pair", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

		pair := make([]models.URLPair, 1)
		pair[0] = testPair
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage.AddBatchURLPairs(context.Background(), pair)
		}
	})

	b.Run("add 100 pairs", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

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
			storage.AddBatchURLPairs(context.Background(), pairs)
		}
	})

	b.Run("add 1000 pairs", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

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
			storage.AddBatchURLPairs(context.Background(), pairs)
		}
	})
}

func TestFileStorage_GetURLPairBatchByUserID(t *testing.T) {
	strg, err := NewFileStorage(testFileName)
	require.NoError(t, err)
	defer os.Remove(strg.strgFileName)
	defer os.Remove(strg.delFileName)

	testPairs := []models.URLPair{
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

	file, err := os.OpenFile(strg.strgFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	require.NoError(t, err)
	defer file.Close()

	enc := json.NewEncoder(file)
	require.NoError(t, err)

	for _, pair := range testPairs {
		err = enc.Encode(&pair)
		require.NoError(t, err)
	}

	t.Run("valid test", func(t *testing.T) {
		pairs, err := strg.GetURLPairBatchByUserID(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, testPairs, pairs)
	})

	t.Run("not exist error", func(t *testing.T) {
		_, err := strg.GetURLPairBatchByUserID(context.Background(), "not exist")
		assert.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("ctx expired", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := strg.GetURLPairBatchByUserID(ctx, testUserID)
		assert.Error(t, err)
	})
}

func BenchmarkFileStorage_GetURLPairBatchByUserID(b *testing.B) {
	b.Run("get url", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

		testPairs := []models.URLPair{
			testPair,
		}

		file, err := os.OpenFile(storage.strgFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		require.NoError(b, err)
		defer file.Close()

		enc := json.NewEncoder(file)
		require.NoError(b, err)

		for _, pair := range testPairs {
			err = enc.Encode(&pair)
			require.NoError(b, err)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage.GetURLPairBatchByUserID(context.Background(), testUserID)
		}
	})

	b.Run("get 100 urls", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

		len := 100
		pairs := make([]models.URLPair, len)
		for i := range len {
			pairs[i] = models.URLPair{
				UID:   models.UserID(fmt.Sprintf("user-%d", i)),
				Orig:  models.OrigURL(fmt.Sprintf("https://site.com/page%d", i)),
				Short: models.ShortURL(fmt.Sprintf("hash-%d", i)),
			}
		}

		file, err := os.OpenFile(storage.strgFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		require.NoError(b, err)
		defer file.Close()

		enc := json.NewEncoder(file)
		require.NoError(b, err)

		for _, pair := range pairs {
			err = enc.Encode(&pair)
			require.NoError(b, err)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage.GetURLPairBatchByUserID(context.Background(), testUserID)
		}
	})

	b.Run("get 1000 urls", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

		len := 1000
		pairs := make([]models.URLPair, len)
		for i := range len {
			pairs[i] = models.URLPair{
				UID:   models.UserID(fmt.Sprintf("user-%d", i)),
				Orig:  models.OrigURL(fmt.Sprintf("https://site.com/page%d", i)),
				Short: models.ShortURL(fmt.Sprintf("hash-%d", i)),
			}
		}

		file, err := os.OpenFile(storage.strgFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		require.NoError(b, err)
		defer file.Close()

		enc := json.NewEncoder(file)
		require.NoError(b, err)

		for _, pair := range pairs {
			err = enc.Encode(&pair)
			require.NoError(b, err)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage.GetURLPairBatchByUserID(context.Background(), testUserID)
		}
	})
}

func TestFileStorage_DeleteRequestedURLs(t *testing.T) {
	strg, err := NewFileStorage(testFileName)
	require.NoError(t, err)
	defer os.Remove(strg.strgFileName)
	defer os.Remove(strg.delFileName)

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

func BenchmarkFileStorage_DeleteRequestedURLs(b *testing.B) {
	b.Run("delete request", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

		delURLs := []*models.DelURLReq{
			&testDelReq,
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage.DeleteRequestedURLs(context.Background(), delURLs)
		}
	})

	b.Run("100 delete requests", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

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
			storage.DeleteRequestedURLs(context.Background(), delURLs)
		}
	})

	b.Run("1000 delete requests", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

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
			storage.DeleteRequestedURLs(context.Background(), delURLs)
		}
	})
}

func TestFileStorage_Ping(t *testing.T) {
	strg, err := NewFileStorage(testFileName)
	require.NoError(t, err)
	defer os.Remove(strg.strgFileName)
	defer os.Remove(strg.delFileName)

	t.Run("valid test", func(t *testing.T) {
		err := strg.Ping(context.Background())
		assert.NoError(t, err)
	})
}

func BenchmarkFileStorage_Ping(b *testing.B) {
	b.Run("ping", func(b *testing.B) {
		storage, err := NewFileStorage(testFileName)
		require.NoError(b, err)
		defer os.Remove(storage.strgFileName)
		defer os.Remove(storage.delFileName)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			storage.Ping(context.Background())
		}
	})
}

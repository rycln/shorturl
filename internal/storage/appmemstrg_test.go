package storage

import (
	"context"
	"testing"

	"github.com/rycln/shorturl/internal/models"
	"github.com/stretchr/testify/assert"
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
		assert.ErrorIs(t, err, ErrConflict)
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
		assert.ErrorIs(t, err, ErrDeletedURL)
	})

	t.Run("not exist error", func(t *testing.T) {
		_, err := strg.GetURLPairByShort(context.Background(), models.ShortURL("not exist"))
		assert.ErrorIs(t, err, ErrNotExist)
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
		assert.ErrorIs(t, err, ErrNotExist)
	})
}

func TestAppMemStorage_DeleteRequestedURLs(t *testing.T) {
	strg := NewAppMemStorage()

	delurls := []models.DelURLReq{
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

func TestAppMemStorage_Ping(t *testing.T) {
	strg := NewAppMemStorage()

	t.Run("valid test", func(t *testing.T) {
		err := strg.Ping(context.Background())
		assert.NoError(t, err)
	})
}

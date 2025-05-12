package storage

import (
	"context"
	"encoding/json"
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

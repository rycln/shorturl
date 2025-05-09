package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/models"
	"github.com/rycln/shorturl/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBatchShortener_BatchShortenURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mHash := mocks.NewMockbatchHasher(ctrl)
	mStrg := mocks.NewMockBatchShortenerStorage(ctrl)

	s := NewBatchShortener(mStrg, mHash)

	testOrigs := []models.OrigURL{
		testOrigURL,
	}

	testPairs := []models.URLPair{
		{
			Short: testShortURL,
			Orig:  testOrigURL,
		},
	}

	t.Run("valid test", func(t *testing.T) {
		mHash.EXPECT().GenerateHashFromURL(testOrigURL).Return(testShortURL)
		mStrg.EXPECT().AddBatchURLPairs(context.Background(), testPairs).Return(nil)

		pairs, err := s.BatchShortenURL(context.Background(), testOrigs)
		assert.NoError(t, err)
		assert.Equal(t, testPairs, pairs)
	})

	t.Run("some error", func(t *testing.T) {
		mHash.EXPECT().GenerateHashFromURL(testOrigURL).Return(testShortURL)
		mStrg.EXPECT().AddBatchURLPairs(context.Background(), testPairs).Return(errTest)

		_, err := s.BatchShortenURL(context.Background(), testOrigs)
		assert.Error(t, err)
	})
}

func TestBatchShortener_GetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mHash := mocks.NewMockbatchHasher(ctrl)
	mStrg := mocks.NewMockBatchShortenerStorage(ctrl)

	s := NewBatchShortener(mStrg, mHash)

	testPairs := []models.URLPair{
		{
			Short: testShortURL,
			Orig:  testOrigURL,
		},
	}

	t.Run("valid test", func(t *testing.T) {
		mStrg.EXPECT().GetURLPairBatchByUserID(context.Background(), testUserID).Return(testPairs, nil)

		pairs, err := s.GetUserURLs(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, testPairs, pairs)
	})

	t.Run("some error", func(t *testing.T) {
		mStrg.EXPECT().GetURLPairBatchByUserID(context.Background(), testUserID).Return(nil, errTest)

		_, err := s.GetUserURLs(context.Background(), testUserID)
		assert.Error(t, err)
	})
}

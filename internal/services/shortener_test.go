package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/shorturl/internal/models"
	"github.com/rycln/shorturl/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShortener_ShortenURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mErr := mocks.NewMockerrConflict(ctrl)
	mHash := mocks.NewMockhasher(ctrl)
	mStrg := mocks.NewMockShortenerStorage(ctrl)

	s := NewShortener(mStrg, mHash)

	wantPair := models.URLPair{
		UID:   testUserID,
		Short: testShortURL,
		Orig:  testOrigURL,
	}

	t.Run("valid test", func(t *testing.T) {
		mHash.EXPECT().GenerateHashFromURL(wantPair.Orig).Return(wantPair.Short)
		mStrg.EXPECT().AddURLPair(context.Background(), &wantPair).Return(nil)

		pair, err := s.ShortenURL(context.Background(), testUserID, wantPair.Orig)
		assert.NoError(t, err)
		assert.Equal(t, &wantPair, pair)
	})

	t.Run("conflict error", func(t *testing.T) {
		mErr.EXPECT().IsErrConflict().Return(true)
		mHash.EXPECT().GenerateHashFromURL(wantPair.Orig).Return(wantPair.Short)
		mStrg.EXPECT().AddURLPair(context.Background(), &wantPair).Return(mErr)

		pair, err := s.ShortenURL(context.Background(), testUserID, wantPair.Orig)
		assert.Error(t, err)
		assert.Equal(t, &wantPair, pair)
	})

	t.Run("some error", func(t *testing.T) {
		mHash.EXPECT().GenerateHashFromURL(wantPair.Orig).Return(wantPair.Short)
		mStrg.EXPECT().AddURLPair(context.Background(), &wantPair).Return(errTest)

		_, err := s.ShortenURL(context.Background(), testUserID, wantPair.Orig)
		assert.Error(t, err)
	})
}

func TestShortener_GetOrigURLByShort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mHash := mocks.NewMockhasher(ctrl)
	mStrg := mocks.NewMockShortenerStorage(ctrl)

	s := NewShortener(mStrg, mHash)

	t.Run("valid test", func(t *testing.T) {
		mStrg.EXPECT().GetURLPairByShort(context.Background(), testShortURL).Return(&testPair, nil)

		orig, err := s.GetOrigURLByShort(context.Background(), testShortURL)
		assert.NoError(t, err)
		assert.Equal(t, testPair.Orig, orig)
	})

	t.Run("some error", func(t *testing.T) {
		mStrg.EXPECT().GetURLPairByShort(context.Background(), testShortURL).Return(nil, errTest)

		_, err := s.GetOrigURLByShort(context.Background(), testShortURL)
		assert.Error(t, err)
	})
}

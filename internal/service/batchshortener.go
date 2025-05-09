package service

import (
	"context"

	"github.com/rycln/shorturl/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type batchURLSaver interface {
	AddBatchURLPairs(context.Context, []models.URLPair) error
}

type batchURLFetcher interface {
	GetURLPairBatchByUserID(context.Context, models.UserID) ([]models.URLPair, error)
}

type BatchShortenerStorage interface {
	batchURLSaver
	batchURLFetcher
}

type batchHasher interface {
	GenerateHashFromURL(models.OrigURL) models.ShortURL
}

type BatchShortener struct {
	strg   BatchShortenerStorage
	hasher batchHasher
}

func NewBatchShortener(strg BatchShortenerStorage, hasher batchHasher) *BatchShortener {
	return &BatchShortener{
		strg:   strg,
		hasher: hasher,
	}
}

func (s *BatchShortener) BatchShortenURL(ctx context.Context, origs []models.OrigURL) ([]models.URLPair, error) {
	var pairs = make([]models.URLPair, len(origs))
	for i, orig := range origs {
		short := s.hasher.GenerateHashFromURL(orig)
		pairs[i] = models.URLPair{
			Short: short,
			Orig:  orig,
		}
	}
	err := s.strg.AddBatchURLPairs(ctx, pairs)
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

func (s *BatchShortener) GetUserURLs(ctx context.Context, uid models.UserID) ([]models.URLPair, error) {
	pairs, err := s.strg.GetURLPairBatchByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

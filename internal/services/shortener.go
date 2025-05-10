package services

import (
	"context"

	"github.com/rycln/shorturl/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type urlSaver interface {
	AddURLPair(context.Context, *models.URLPair) error
}

type urlFetcher interface {
	GetURLPairByShort(context.Context, models.ShortURL) (*models.URLPair, error)
}

type ShortenerStorage interface {
	urlSaver
	urlFetcher
}

type hasher interface {
	GenerateHashFromURL(models.OrigURL) models.ShortURL
}

type errConflict interface {
	error
	IsConflict() bool
}

type Shortener struct {
	strg   ShortenerStorage
	hasher hasher
}

func NewShortener(strg ShortenerStorage, hasher hasher) *Shortener {
	return &Shortener{
		strg:   strg,
		hasher: hasher,
	}
}

func (s *Shortener) ShortenURL(ctx context.Context, orig models.OrigURL) (*models.URLPair, error) {
	short := s.hasher.GenerateHashFromURL(orig)
	pair := &models.URLPair{
		Short: short,
		Orig:  orig,
	}
	err := s.strg.AddURLPair(ctx, pair)
	if e, ok := err.(errConflict); ok && e.IsConflict() {
		return pair, err
	}
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func (s *Shortener) GetOrigURLByShort(ctx context.Context, short models.ShortURL) (models.OrigURL, error) {
	pair, err := s.strg.GetURLPairByShort(ctx, short)
	if err != nil {
		return "", err
	}
	return pair.Orig, nil
}

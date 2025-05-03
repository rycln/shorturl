package service

import (
	"context"

	"github.com/rycln/shorturl/internal/models"
)

type URLSaver interface {
	AddURLPair(context.Context, *models.URLPair) error
}

type URLFetcher interface {
	GetURLPairByShort(context.Context, models.ShortURL) (*models.URLPair, error)
}

type Hasher interface {
	GenerateHashFromURL(models.OrigURL) models.ShortURL
}

type errConflict interface {
	error
	IsConflict() bool
}

type Shortener struct {
	saver   URLSaver
	fetcher URLFetcher
	hasher  Hasher
}

func NewShortener(saver URLSaver, fetcher URLFetcher) *Shortener {
	return &Shortener{
		saver:   saver,
		fetcher: fetcher,
	}
}

func (s *Shortener) ShortenURL(ctx context.Context, orig models.OrigURL) (*models.URLPair, error) {
	short := s.hasher.GenerateHashFromURL(orig)
	pair := &models.URLPair{
		Short: short,
		Orig:  orig,
	}
	err := s.saver.AddURLPair(ctx, pair)
	if e, ok := err.(errConflict); ok && e.IsConflict() {
		return pair, err
	}
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func (s *Shortener) GetOrigURLByShort(ctx context.Context, short models.ShortURL) (models.OrigURL, error) {
	pair, err := s.fetcher.GetURLPairByShort(ctx, short)
	if err != nil {
		return "", err
	}
	return pair.Orig, nil
}

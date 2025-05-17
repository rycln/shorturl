package services

import (
	"context"
	"errors"

	"github.com/rycln/shorturl/internal/contextkeys"
	"github.com/rycln/shorturl/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

var ErrNoShortURL = errors.New("short URL value is empty")

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
	IsErrConflict() bool
}

// Shortener provides core URL shortening functionality.
//
// The service handles all business logic for URL operations including:
// - Creating shortened URLs from original ones
// - Retrieving original URLs from shortened versions
// - Context-based URL resolution for HTTP handlers
type Shortener struct {
	strg   ShortenerStorage
	hasher hasher
}

// NewShortener creates a new Shortener service instance.
func NewShortener(strg ShortenerStorage, hasher hasher) *Shortener {
	return &Shortener{
		strg:   strg,
		hasher: hasher,
	}
}

// Shorten creates a URLpair instance from original URL and user id.
//
// Returns the shortened URL pair or error if operation fails.
func (s *Shortener) ShortenURL(ctx context.Context, uid models.UserID, orig models.OrigURL) (*models.URLPair, error) {
	short := s.hasher.GenerateHashFromURL(orig)
	pair := &models.URLPair{
		UID:   uid,
		Short: short,
		Orig:  orig,
	}
	err := s.strg.AddURLPair(ctx, pair)
	if e, ok := err.(errConflict); ok && e.IsErrConflict() {
		return pair, err
	}
	if err != nil {
		return nil, err
	}
	return pair, nil
}

// GetOrigURLByShort retrieves the original URL from a shortened version.
//
// Returns the original URL if found in storage.
// Returns error if short URL is invalid or not found.
func (s *Shortener) GetOrigURLByShort(ctx context.Context, short models.ShortURL) (models.OrigURL, error) {
	pair, err := s.strg.GetURLPairByShort(ctx, short)
	if err != nil {
		return "", err
	}
	return pair.Orig, nil
}

// GetShortURLFromCtx extracts shortened URL from request context.
//
// Returns empty string and error if URL not found in context.
//
// Typical usage:
//
//	shortURL, err := s.GetShortURLFromContext(r.Context())
func (s *Shortener) GetShortURLFromCtx(ctx context.Context) (models.ShortURL, error) {
	shortURL, ok := ctx.Value(contextkeys.ShortURL).(string)
	if !ok {
		return "", ErrNoShortURL
	}
	return models.ShortURL(shortURL), nil
}

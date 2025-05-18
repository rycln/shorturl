package services

import (
	"context"

	"github.com/rycln/shorturl/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// batchURLSaver defines batch URL storage operations.
type batchURLSaver interface {
	// AddBatchURLPairs stores multiple URL pairs in single transaction.
	AddBatchURLPairs(context.Context, []models.URLPair) error
}

// batchURLFetcher defines user URL retrieval operations.
type batchURLFetcher interface {
	// GetURLPairBatchByUserID retrieves all shortened URLs for a specific user.
	GetURLPairBatchByUserID(context.Context, models.UserID) ([]models.URLPair, error)
}

// BatchShortenerStorage combines storage operations needed for batch URL processing.
//
// The interface composes two fundamental capabilities required by the BatchShortener service:
//   - Saving multiple URLs in single operation (batchURLSaver)
//   - Retrieving user's URLs (batchURLFetcher)
type BatchShortenerStorage interface {
	batchURLSaver
	batchURLFetcher
}

type batchHasher interface {
	GenerateHashFromURL(models.OrigURL) models.ShortURL
}

// BatchShortener provides batch operations for URL shortening.
//
// The service handles processing of multiple URLs in single operation,
// optimizing storage access and maintaining consistency.
type BatchShortener struct {
	strg   BatchShortenerStorage
	hasher batchHasher
}

// NewBatchShortener creates new batch processor instance.
func NewBatchShortener(strg BatchShortenerStorage, hasher batchHasher) *BatchShortener {
	return &BatchShortener{
		strg:   strg,
		hasher: hasher,
	}
}

// BatchShortenURL processes multiple URLs in single operation.
//
// Accepts slice of original URLs and user ID that owns them.
// Returns slice of URLPair structures containing both original
// and shortened versions, maintaining input order.
func (s *BatchShortener) BatchShortenURL(ctx context.Context, uid models.UserID, origs []models.OrigURL) ([]models.URLPair, error) {
	var pairs = make([]models.URLPair, len(origs))
	for i, orig := range origs {
		short := s.hasher.GenerateHashFromURL(orig)
		pairs[i] = models.URLPair{
			UID:   uid,
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

// GetUserURLs retrieves all shortened URLs for specific user.
//
// Returns slice of URLPair structures or empty slice if none found.
func (s *BatchShortener) GetUserURLs(ctx context.Context, uid models.UserID) ([]models.URLPair, error) {
	pairs, err := s.strg.GetURLPairBatchByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

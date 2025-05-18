package services

import (
	"context"

	"github.com/rycln/shorturl/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// BatchDeleterStorage defines the storage interface required by BatchDeleter service.
type BatchDeleterStorage interface {
	// DeleteURLs marks multiple URLs as deleted in storage.
	DeleteRequestedURLs(context.Context, []*models.DelURLReq) error
}

// BatchDeleter provides batch deletion functionality for shortened URLs.
//
// The service handles asynchronous deletion of multiple URLs.
type BatchDeleter struct {
	strg BatchDeleterStorage
}

// NewBatchDeleter creates new batch deletion service instance.
func NewBatchDeleter(strg BatchDeleterStorage) *BatchDeleter {
	return &BatchDeleter{
		strg: strg,
	}
}

// DeleteBatch processes multiple URL deletion requests.
//
// Accepts slice of DelURLReq structures containing:
// - ShortURL to delete
// - UserID for ownership verification
func (s *BatchDeleter) DeleteURLsBatch(ctx context.Context, urls []*models.DelURLReq) error {
	err := s.strg.DeleteRequestedURLs(ctx, urls)
	if err != nil {
		return err
	}
	return nil
}

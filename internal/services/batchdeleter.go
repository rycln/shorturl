package services

import (
	"context"

	"github.com/rycln/shorturl/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type BatchDeleterStorage interface {
	DeleteRequestedURLs(context.Context, []*models.DelURLReq) error
}

type BatchDeleter struct {
	strg BatchDeleterStorage
}

func NewBatchDeleter(strg BatchDeleterStorage) *BatchDeleter {
	return &BatchDeleter{
		strg: strg,
	}
}

func (s *BatchDeleter) DeleteURLsBatch(ctx context.Context, urls []*models.DelURLReq) error {
	err := s.strg.DeleteRequestedURLs(ctx, urls)
	if err != nil {
		return err
	}
	return nil
}

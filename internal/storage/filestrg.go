package storage

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/rycln/shorturl/internal/models"
)

// FileStorage is a persistent file-based implementation of a URL shortener storage.
// It provides operations for storing and retrieving URL pairs with disk persistence.
type FileStorage struct {
	strgFileName string
	delFileName  string
	strgMu       sync.Mutex
	delMu        sync.Mutex
}

// NewFileStorage creates a new FileStorage instance.
func NewFileStorage(fileName string) (*FileStorage, error) {
	_, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	delFileName := fileName + "_deleted"

	_, err = os.Create(delFileName)
	if err != nil {
		return nil, err
	}

	return &FileStorage{
		strgFileName: fileName,
		delFileName:  delFileName,
	}, nil
}

// AddURLPair stores a new URL pair in the file storage.
func (s *FileStorage) AddURLPair(ctx context.Context, pair *models.URLPair) error {
	_, err := s.getPairByShort(ctx, pair.Short)
	if errors.Is(err, errNotExist) {
		err = s.writeIntoStrgFile(pair)
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	return newErrConflict(errConflict)
}

// GetURLPairByShort retrieves a URL pair by its short URL from file storage.
func (s *FileStorage) GetURLPairByShort(ctx context.Context, short models.ShortURL) (*models.URLPair, error) {
	deleted, err := s.shortIsDeleted(ctx, short)
	if err != nil {
		return nil, err
	}
	if deleted {
		return nil, newErrDeletedURL(errDeletedURL)
	}

	pair, err := s.getPairByShort(ctx, short)
	if err != nil {
		return nil, err
	}
	return pair, nil
}

// AddBatchURLPairs stores multiple URL pairs in a single file operation.
func (s *FileStorage) AddBatchURLPairs(ctx context.Context, pairs []models.URLPair) error {
	for _, pair := range pairs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := s.writeIntoStrgFile(&pair)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetURLPairBatchByUserID retrieves all URL pairs created by a specific user.
func (s *FileStorage) GetURLPairBatchByUserID(ctx context.Context, uid models.UserID) ([]models.URLPair, error) {
	return s.getAllUserPairs(ctx, uid)
}

// DeleteRequestedURLs marks URLs as deleted in a batch operation.
// Implements soft deletion - URLs remain in storage but are marked as deleted.
func (s *FileStorage) DeleteRequestedURLs(ctx context.Context, delurls []*models.DelURLReq) error {
	for _, delurl := range delurls {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := s.writeIntoDelFile(delurl)
		if err != nil {
			return err
		}
	}
	return nil
}

// Ping is a no-op health check that always succeeds for file storage.
// Exists to satisfy storage interface requirements.
func (s *FileStorage) Ping(context.Context) error { return nil }

// Close is a no-op cleanup method for file storage.
// Exists to satisfy storage interface requirements.
func (s *FileStorage) Close() error { return nil }

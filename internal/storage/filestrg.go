package storage

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/rycln/shorturl/internal/models"
)

type FileStorage struct {
	strgMu       sync.Mutex
	delMu        sync.Mutex
	strgFileName string
	delFileName  string
}

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

func (s *FileStorage) AddURLPair(ctx context.Context, pair *models.URLPair) error {
	_, err := s.getPairByShort(ctx, pair.Short)
	if errors.Is(err, ErrNotExist) {
		err := s.writeIntoStrgFile(pair)
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	return newErrConflict(ErrConflict)
}

func (s *FileStorage) GetURLPairByShort(ctx context.Context, short models.ShortURL) (*models.URLPair, error) {
	deleted, err := s.shortIsDeleted(ctx, short)
	if err != nil {
		return nil, err
	}
	if deleted {
		return nil, newErrDeletedURL(ErrDeletedURL)
	}

	pair, err := s.getPairByShort(ctx, short)
	if err != nil {
		return nil, err
	}
	return pair, nil
}

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

func (s *FileStorage) GetURLPairBatchByUserID(ctx context.Context, uid models.UserID) ([]models.URLPair, error) {
	return s.getAllUserPairs(ctx, uid)
}

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

func (s *FileStorage) Ping(context.Context) error { return nil }

func (s *FileStorage) Close() {}

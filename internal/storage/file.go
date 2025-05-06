package storage

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/rycln/shorturl/internal/models"
)

type FileStorage struct {
	mu         sync.Mutex
	enc        *fileEncoder
	decFactory *fileDecoderFactory
}

func NewFileStorage(fileName string) *FileStorage {
	enc, err := newFileEncoder(fileName)
	if err != nil {
		log.Fatalf("Can't open file: %v", err)
	}

	decFactory := newFileDecoderFactory(fileName)

	return &FileStorage{
		enc:        enc,
		decFactory: decFactory,
	}
}

func (s *FileStorage) AddURLPair(ctx context.Context, pair *models.URLPair) error {
	_, err := s.getPairByOrig(ctx, pair.Orig)
	if errors.Is(err, ErrNotExist) {
		err := s.writeIntoFile(pair)
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
	pair, err := s.getPairByShort(ctx, short)
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func (s *FileStorage) AddBatchURLPairs(ctx context.Context, pairs []models.URLPair) error {
	for _, pair := range pairs {
		err := s.writeIntoFile(&pair)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FileStorage) GetURLPairBatchByUserID(ctx context.Context, uid models.UserID) ([]models.URLPair, error) {
	return s.getAllUserPairs(ctx, uid)
}

func (s *FileStorage) DeleteRequestedURLs(ctx context.Context, delurls []models.DelURLReq) error {

}

func (s *FileStorage) Ping(context.Context) error { return nil }

func (s *FileStorage) Close() {
	s.enc.close()
}

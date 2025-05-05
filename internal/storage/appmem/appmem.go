package appmem

import (
	"context"
	"sync"

	"github.com/rycln/shorturl/internal/models"
)

type AppMemStorage struct {
	mu          sync.RWMutex
	globalPairs map[models.ShortURL]models.OrigURL
	userPairs   map[models.UserID]map[models.ShortURL]models.OrigURL
}

func NewAppMemStorage() *AppMemStorage {
	return &AppMemStorage{
		globalPairs: make(map[models.ShortURL]models.OrigURL),
		userPairs:   make(map[models.UserID]map[models.ShortURL]models.OrigURL),
	}
}

func (s *AppMemStorage) AddURLPair(context.Context, *models.URLPair) error {
	_, ok := ss.storage[surl.ShortURL]
	if ok {
		return ErrConflict
	}
	ss.storage[surl.ShortURL] = surl.OrigURL
	return nil
}

func (s *AppMemStorage) AddBatchURL(ctx context.Context, surls []ShortenedURL) error {
	for _, surl := range surls {
		ss.storage[surl.ShortURL] = surl.OrigURL
	}
	return nil
}

func (s *AppMemStorage) GetOrigURL(ctx context.Context, shortURL string) (string, error) {
	_, ok := ss.storage[shortURL]
	if !ok {
		return "", ErrNotExist
	}
	origURL := ss.storage[shortURL]
	return origURL, nil
}

func (s *AppMemStorage) GetShortURL(ctx context.Context, origURL string) (string, error) {
	for short, orig := range ss.storage {
		if orig == origURL {
			return short, nil
		}
		select {
		case <-ctx.Done():
			return "", ErrTimeLimit
		default:
			continue
		}
	}
	return "", ErrNotExist
}

func (s *AppMemStorage) Ping() {

}

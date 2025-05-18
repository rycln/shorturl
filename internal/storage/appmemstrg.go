package storage

import (
	"context"
	"sync"

	"github.com/rycln/shorturl/internal/models"
)

// AppMemStorage is an in-memory implementation of a URL shortener storage.
//
// Note: All data will be lost on application restart.
type AppMemStorage struct {
	mu      sync.RWMutex
	pairs   map[models.UserID]map[models.ShortURL]models.OrigURL
	deleted map[models.ShortURL]struct{}
}

// NewAppMemStorage creates a new AppMemStorage instance.
func NewAppMemStorage() *AppMemStorage {
	return &AppMemStorage{
		pairs:   make(map[models.UserID]map[models.ShortURL]models.OrigURL),
		deleted: make(map[models.ShortURL]struct{}),
	}
}

// AddURLPair stores a new URL pair in memory.
func (s *AppMemStorage) AddURLPair(ctx context.Context, pair *models.URLPair) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if userpairs, exists := s.pairs[pair.UID]; exists {
		if _, ok := userpairs[pair.Short]; ok {
			return newErrConflict(errConflict)
		}

		userpairs[pair.Short] = pair.Orig
		return nil
	}

	for _, userpairs := range s.pairs {
		if _, ok := userpairs[pair.Short]; ok {
			return newErrConflict(errConflict)
		}
	}

	s.pairs[pair.UID] = map[models.ShortURL]models.OrigURL{
		pair.Short: pair.Orig,
	}

	return nil
}

// GetURLPairByShort retrieves a URL pair by its short URL.
func (s *AppMemStorage) GetURLPairByShort(ctx context.Context, short models.ShortURL) (*models.URLPair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.deleted[short]
	if ok {
		return nil, newErrDeletedURL(errDeletedURL)
	}

	for uid, userpairs := range s.pairs {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		_, ok := userpairs[short]
		if ok {
			var pair = &models.URLPair{
				UID:   uid,
				Short: short,
				Orig:  userpairs[short],
			}
			return pair, nil
		}
	}

	return nil, newErrNotExist(errNotExist)
}

// AddBatchURLPairs stores multiple URL pairs.
func (s *AppMemStorage) AddBatchURLPairs(ctx context.Context, pairs []models.URLPair) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, pair := range pairs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, ok := s.pairs[pair.UID]
		if ok {
			s.pairs[pair.UID][pair.Short] = pair.Orig
			continue
		}

		var userpair = make(map[models.ShortURL]models.OrigURL)
		userpair[pair.Short] = pair.Orig
		s.pairs[pair.UID] = userpair
	}

	return nil
}

// GetURLPairBatchByUserID retrieves all URL pairs created by a specific user.
func (s *AppMemStorage) GetURLPairBatchByUserID(ctx context.Context, uid models.UserID) ([]models.URLPair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.pairs[uid]
	if !ok {
		return nil, newErrNotExist(errNotExist)
	}

	var pairs []models.URLPair

	for short, orig := range s.pairs[uid] {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		pair := models.URLPair{
			UID:   uid,
			Short: short,
			Orig:  orig,
		}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

// DeleteRequestedURLs marks URLs as deleted in a batch operation.
// Implements soft deletion - URLs remain in storage but are marked as deleted.
func (s *AppMemStorage) DeleteRequestedURLs(ctx context.Context, delurls []*models.DelURLReq) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, delurl := range delurls {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		s.deleted[delurl.Short] = struct{}{}
	}

	return nil
}

// Ping is a no-op health check that always succeeds for in-memory storage.
// Exists to satisfy storage interface requirements.
func (s *AppMemStorage) Ping(context.Context) error { return nil }

// Close is a no-op cleanup method for in-memory storage.
// Exists to satisfy storage interface requirements.
func (s *AppMemStorage) Close() {}

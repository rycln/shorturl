package storage

import (
	"context"
	"sync"

	"github.com/rycln/shorturl/internal/models"
)

type AppMemStorage struct {
	mu    sync.RWMutex
	pairs map[models.UserID]map[models.ShortURL]models.OrigURL
}

func NewAppMemStorage() *AppMemStorage {
	return &AppMemStorage{
		pairs: make(map[models.UserID]map[models.ShortURL]models.OrigURL),
	}
}

func (s *AppMemStorage) AddURLPair(ctx context.Context, pair *models.URLPair) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, userpairs := range s.pairs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, ok := userpairs[pair.Short]
		if ok {
			return newErrConflict(ErrConflict)
		}
	}

	_, ok := s.pairs[pair.UID]
	if ok {
		s.pairs[pair.UID][pair.Short] = pair.Orig
		return nil
	}

	var userpair = make(map[models.ShortURL]models.OrigURL)
	userpair[pair.Short] = pair.Orig
	s.pairs[pair.UID] = userpair

	return nil
}

func (s *AppMemStorage) GetURLPairByShort(ctx context.Context, short models.ShortURL) (*models.URLPair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

	return nil, newErrNotExist(ErrNotExist)
}

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

func (s *AppMemStorage) GetURLPairBatchByUserID(ctx context.Context, uid models.UserID) ([]models.URLPair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.pairs[uid]
	if !ok {
		return nil, newErrNotExist(ErrNotExist)
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

func (s *AppMemStorage) DeleteRequestedURLs(ctx context.Context, delurls []models.DelURLReq) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, delurl := range delurls {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		delete(s.pairs[delurl.UID], delurl.Short)
	}

	return nil
}

func (s *AppMemStorage) Ping(context.Context) error { return nil }

func (s *AppMemStorage) Close() {}

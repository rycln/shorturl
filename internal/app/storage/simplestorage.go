package storage

import (
	"context"
)

type SimpleStorage struct {
	storage map[string]string
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		storage: make(map[string]string),
	}
}

func (ss SimpleStorage) AddURL(ctx context.Context, surl ShortenedURL) error {
	_, ok := ss.storage[surl.ShortURL]
	if ok {
		return ErrConflict
	}
	ss.storage[surl.ShortURL] = surl.OrigURL
	return nil
}

func (ss SimpleStorage) AddBatchURL(ctx context.Context, surls []ShortenedURL) error {
	for _, surl := range surls {
		ss.storage[surl.ShortURL] = surl.OrigURL
	}
	return nil
}

func (ss SimpleStorage) GetOrigURL(ctx context.Context, shortURL string) (string, error) {
	_, ok := ss.storage[shortURL]
	if !ok {
		return "", ErrNotExist
	}
	origURL := ss.storage[shortURL]
	return origURL, nil
}

func (ss SimpleStorage) GetShortURL(ctx context.Context, origURL string) (string, error) {
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

func (ss SimpleStorage) Ping(ctx context.Context) error {
	return ErrNotDatabase
}

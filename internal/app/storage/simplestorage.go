package storage

import (
	"context"
	"errors"
)

type SimpleStorage struct {
	storage map[string]string
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		storage: make(map[string]string),
	}
}

func (ss SimpleStorage) AddURL(ctx context.Context, surls ...ShortenedURL) error {
	for _, surl := range surls {
		_, ok := ss.storage[surl.ShortURL]
		if ok {
			return errors.New("shortened URL already exist")
		}
		ss.storage[surl.ShortURL] = surl.OrigURL
	}
	return nil
}

func (ss SimpleStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	_, ok := ss.storage[shortURL]
	if !ok {
		return "", errors.New("shortened URL does not exist")
	}
	origURL := ss.storage[shortURL]
	return origURL, nil
}

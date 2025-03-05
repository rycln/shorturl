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

func (sms SimpleStorage) AddURL(ctx context.Context, shortURL, fullURL string) error {
	_, ok := sms.storage[shortURL]
	if ok {
		return errors.New("shortened URL already exist")
	}
	sms.storage[shortURL] = fullURL
	return nil
}

func (sms SimpleStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	_, ok := sms.storage[shortURL]
	if !ok {
		return "", errors.New("shortened URL does not exist")
	}
	return sms.storage[shortURL], nil
}

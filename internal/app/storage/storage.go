package storage

import (
	"database/sql"
	"errors"
)

var DB *sql.DB

type SimpleMemStorage struct {
	storage map[string]string
}

func NewSimpleMemStorage() *SimpleMemStorage {
	return &SimpleMemStorage{
		storage: make(map[string]string),
	}
}

func (sms SimpleMemStorage) AddURL(shortURL, fullURL string) bool {
	_, ok := sms.storage[shortURL]
	if ok {
		return false
	}
	sms.storage[shortURL] = fullURL
	return true
}

func (sms SimpleMemStorage) GetURL(shortURL string) (string, error) {
	_, ok := sms.storage[shortURL]
	if !ok {
		return "", errors.New("shortened URL does not exist")
	}
	return sms.storage[shortURL], nil
}

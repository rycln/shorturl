package mem

import (
	"errors"
)

type Storager interface {
	AddURL(string, string)
	GetURL(string) (string, error)
}

type SimpleMemStorage struct {
	storage map[string]string
}

func (sms SimpleMemStorage) AddURL(shortURL, fullURL string) {
	_, ok := sms.storage[shortURL]
	if ok {
		return
	}
	sms.storage[shortURL] = fullURL
}

func (sms SimpleMemStorage) GetURL(shortURL string) (string, error) {
	_, ok := sms.storage[shortURL]
	if !ok {
		return "", errors.New("shortened URL does not exist")
	}
	return sms.storage[shortURL], nil
}

func NewSimpleMemStorage() SimpleMemStorage {
	return SimpleMemStorage{
		storage: make(map[string]string),
	}
}

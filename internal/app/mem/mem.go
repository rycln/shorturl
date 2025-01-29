package mem

import (
	"errors"

	"github.com/rycln/shorturl/internal/app/hash"
)

type Memorizer interface {
	AddURL(string)
	GetURL(string) string
}

type MemStorage struct {
	Storage map[string]string
	incId   int64
}

func (m MemStorage) AddURL(url string) string {
	shortURL := hash.Base62(m.incId)
	m.Storage[shortURL] = url
	return shortURL
}

func (m MemStorage) GetURL(shortURL string) (string, error) {
	_, ok := m.Storage[shortURL]
	if !ok {
		return "", errors.New("wrong shortened URL")
	}
	return m.Storage[shortURL], nil
}

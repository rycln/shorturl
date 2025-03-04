package storage

import (
	"github.com/google/uuid"
)

type StoredURL struct {
	ID       string `json:"uuid"`
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}

func NewStoredURL(shortURL, fullURL string) *StoredURL {
	surl := &StoredURL{
		ID:       uuid.NewString(),
		ShortURL: shortURL,
		FullURL:  fullURL,
	}
	return surl
}

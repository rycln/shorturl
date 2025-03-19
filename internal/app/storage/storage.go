package storage

import (
	"errors"
)

var (
	ErrConflict  = errors.New("shortened URL already exists")
	ErrNotExist  = errors.New("shortened URL does not exist")
	ErrTimeLimit = errors.New("time limit exceeded")
)

type ShortenedURL struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"original_url"`
}

func NewShortenedURL(shortURL, origURL string) ShortenedURL {
	surl := ShortenedURL{
		ShortURL: shortURL,
		OrigURL:  origURL,
	}
	return surl
}

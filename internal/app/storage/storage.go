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
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"original_url"`
}

func NewShortenedURL(uid, shortURL, origURL string) ShortenedURL {
	surl := ShortenedURL{
		UserID:   uid,
		ShortURL: shortURL,
		OrigURL:  origURL,
	}
	return surl
}

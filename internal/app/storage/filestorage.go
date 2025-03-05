package storage

import (
	"context"

	"github.com/google/uuid"
)

type storedURL struct {
	ID       string `json:"uuid"`
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}

func newStoredURL(shortURL, fullURL string) *storedURL {
	surl := &storedURL{
		ID:       uuid.NewString(),
		ShortURL: shortURL,
		FullURL:  fullURL,
	}
	return surl
}

type FileStorage struct {
	encoder *FileEncoder
	decoder *FileDecoder
}

func NewFileStorage(enc *FileEncoder, dec *FileDecoder) *FileStorage {
	return &FileStorage{
		encoder: enc,
		decoder: dec,
	}
}

func (fs *FileStorage) AddURL(ctx context.Context, shortURL, fullURL string) error {
	surl := newStoredURL(shortURL, fullURL)
	err := fs.encoder.writeIntoFile(surl)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	return fs.decoder.getFromFile(ctx, shortURL)
}

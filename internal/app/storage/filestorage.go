package storage

import (
	"context"
)

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

func (fs *FileStorage) AddURL(ctx context.Context, surls ...ShortenedURL) error {
	for _, surl := range surls {
		err := fs.encoder.writeIntoFile(&surl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	return fs.decoder.getFromFile(ctx, shortURL)
}

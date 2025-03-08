package storage

import (
	"context"
	"errors"
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

func (fs *FileStorage) AddURL(ctx context.Context, surl ShortenedURL) error {
	checkURL, err := fs.decoder.getFromFile(ctx, surl.OrigURL)
	if err != nil {
		if !errors.Is(err, ErrNotExist) {
			return err
		}
	}
	if checkURL != nil {
		return ErrConflict
	}
	err = fs.encoder.writeIntoFile(&surl)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) AddBatchURL(ctx context.Context, surls []ShortenedURL) error {
	for _, surl := range surls {
		err := fs.encoder.writeIntoFile(&surl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileStorage) GetOrigURL(ctx context.Context, shortURL string) (string, error) {
	surl, err := fs.decoder.getFromFile(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return surl.OrigURL, nil
}

func (fs *FileStorage) GetShortURL(ctx context.Context, origURL string) (string, error) {
	surl, err := fs.decoder.getFromFile(ctx, origURL)
	if err != nil {
		return "", err
	}
	return surl.ShortURL, nil
}

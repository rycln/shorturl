package storage

import (
	"context"
	"errors"
	"io"
	"log"
)

type FileStorage struct {
	fileName string
	encoder  *fileEncoder
}

func NewFileStorage(fileName string) (*FileStorage, func() error) {
	encoder, err := newFileEncoder(fileName)
	if err != nil {
		log.Fatalf("Can't open file: %v", err)
	}
	return &FileStorage{
		fileName: fileName,
		encoder:  encoder,
	}, encoder.close
}

func (fs *FileStorage) AddURL(ctx context.Context, surl ShortenedURL) error {
	_, err := fs.getFromFile(ctx, surl.OrigURL)
	if err != nil {
		if errors.Is(err, ErrNotExist) {
			err := fs.encoder.writeIntoFile(&surl)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return ErrConflict
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
	surl, err := fs.getFromFile(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return surl.OrigURL, nil
}

func (fs *FileStorage) GetShortURL(ctx context.Context, origURL string) (string, error) {
	surl, err := fs.getFromFile(ctx, origURL)
	if err != nil {
		return "", err
	}
	return surl.ShortURL, nil
}

func (fs *FileStorage) getFromFile(ctx context.Context, url string) (*ShortenedURL, error) {
	fd, err := newFileDecoder(fs.fileName)
	if err != nil {
		return nil, err
	}
	defer fd.close()

	for {
		surl := &ShortenedURL{}
		err := fd.decoder.Decode(surl)
		if err == io.EOF {
			return nil, ErrNotExist
		}
		if err != nil {
			return nil, err
		}
		if surl.ShortURL == url || surl.OrigURL == url {
			return surl, nil
		}
		select {
		case <-ctx.Done():
			return nil, ErrTimeLimit
		default:
			continue
		}
	}
}

func (fs *FileStorage) GetAllUserURLs(ctx context.Context, uid string) ([]ShortenedURL, error) {
	fd, err := newFileDecoder(fs.fileName)
	if err != nil {
		return nil, err
	}
	defer fd.close()

	var surls []ShortenedURL
	for {
		surl := &ShortenedURL{}
		err := fd.decoder.Decode(surl)
		if err == io.EOF {
			if surls == nil {
				return nil, ErrNotExist
			}
			return surls, nil
		}
		if err != nil {
			return nil, err
		}
		if surl.UserID == uid {
			surls = append(surls, *surl)
		}
		select {
		case <-ctx.Done():
			return nil, ErrTimeLimit
		default:
			continue
		}
	}
}

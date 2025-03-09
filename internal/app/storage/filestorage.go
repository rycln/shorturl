package storage

import (
	"context"
	"errors"
	"io"
)

type FileStorage struct {
	fileName string
	encoder  *fileEncoder
}

func NewFileStorage(fileName string) (*FileStorage, error) {
	encoder, err := newFileEncoder(fileName)
	if err != nil {
		return nil, err
	}
	return &FileStorage{
		fileName: fileName,
		encoder:  encoder,
	}, nil
}

func (fs *FileStorage) Close() {
	fs.encoder.close()
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
		if err != nil {
			if err == io.EOF {
				return nil, ErrNotExist
			}
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

package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type FileDecoder struct {
	file    *os.File
	decoder *json.Decoder
}

func NewFileDecoder(fileName string) (*FileDecoder, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileDecoder{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (fd *FileDecoder) Close() error {
	return fd.file.Close()
}

func (fd *FileDecoder) getFromFile(ctx context.Context, url string) (*ShortenedURL, error) {
	for {
		surl := &ShortenedURL{}
		err := fd.decoder.Decode(surl)
		if err != nil {
			if errors.Is(err, io.EOF) {
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

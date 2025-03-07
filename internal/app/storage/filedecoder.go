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

func (fd *FileDecoder) getFromFile(ctx context.Context, shortURL string) (string, error) {
	for {
		surl := &ShortenedURL{}
		err := fd.decoder.Decode(surl)
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			return "", errors.New("shortened URL does not exist")
		}
		if surl.ShortURL == shortURL {
			return surl.OrigURL, nil
		}
		select {
		case <-ctx.Done():
			return "", errors.New("time limit exceeded")
		default:
			continue
		}
	}
}

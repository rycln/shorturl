package storage

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/rycln/shorturl/internal/models"
)

type fileDecoderFactory struct {
	file *os.File
}

func newFileDecoderFactory(fileName string) (*fileDecoderFactory, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return &fileDecoderFactory{file: file}, nil
}

func (f *fileDecoderFactory) NewDecoder() (*json.Decoder, error) {
	_, err := f.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return json.NewDecoder(f.file), nil
}

func (fd *fileDecoder) close() error {
	return fd.file.Close()
}

func getPairFromFileWithCondition(ctx context.Context, fileName string, condition func(*models.URLPair) bool) (*models.URLPair, error) {
	fd, err := newFileDecoder(fileName)
	if err != nil {
		return nil, err
	}
	defer fd.close()

	for {
		pair := &models.URLPair{}
		err := fd.decoder.Decode(pair)
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

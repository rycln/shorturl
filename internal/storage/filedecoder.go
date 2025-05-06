package storage

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/rycln/shorturl/internal/models"
)

type fileDecoder struct {
	*json.Decoder
	file *os.File
}

type fileDecoderFactory struct {
	fileName string
}

func newFileDecoderFactory(fileName string) *fileDecoderFactory {
	return &fileDecoderFactory{fileName: fileName}
}

func (f *fileDecoderFactory) newFileDecoder() (*fileDecoder, error) {
	file, err := os.Open(f.fileName)
	if err != nil {
		return nil, err
	}

	return &fileDecoder{
		file:    file,
		Decoder: json.NewDecoder(file),
	}, nil
}

func (f *fileDecoder) close() error {
	return f.file.Close()
}

func (s *FileStorage) getPairByShort(ctx context.Context, short models.ShortURL) (*models.URLPair, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fd, err := s.decFactory.newFileDecoder()
	if err != nil {
		return nil, err
	}
	defer fd.close()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		pair := &models.URLPair{}
		err := fd.Decode(pair)
		if err == io.EOF {
			return nil, ErrNotExist
		}
		if err != nil {
			return nil, err
		}

		if pair.Short == short {
			return pair, nil
		}
	}
}

func (s *FileStorage) getPairByOrig(ctx context.Context, orig models.OrigURL) (*models.URLPair, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fd, err := s.decFactory.newFileDecoder()
	if err != nil {
		return nil, err
	}
	defer fd.close()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		pair := &models.URLPair{}
		err := fd.Decode(pair)
		if err == io.EOF {
			return nil, ErrNotExist
		}
		if err != nil {
			return nil, err
		}

		if pair.Orig == orig {
			return pair, nil
		}
	}
}

func (s *FileStorage) getAllUserPairs(ctx context.Context, uid models.UserID) ([]models.URLPair, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fd, err := s.decFactory.newFileDecoder()
	if err != nil {
		return nil, err
	}
	defer fd.close()

	var userpairs []models.URLPair

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		pair := &models.URLPair{}
		err := fd.Decode(pair)
		if err == io.EOF {
			if userpairs == nil {
				return nil, ErrNotExist
			}
			return userpairs, nil
		}
		if err != nil {
			return nil, err
		}

		if pair.UID == uid {
			userpairs = append(userpairs, *pair)
		}
	}
}

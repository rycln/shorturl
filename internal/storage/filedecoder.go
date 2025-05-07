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

func newFileDecoder(fileName string) (*fileDecoder, error) {
	file, err := os.Open(fileName)
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
	s.strgMu.Lock()
	defer s.strgMu.Unlock()

	fd, err := newFileDecoder(s.strgFileName)
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

func (s *FileStorage) getAllUserPairs(ctx context.Context, uid models.UserID) ([]models.URLPair, error) {
	s.strgMu.Lock()
	defer s.strgMu.Unlock()

	fd, err := newFileDecoder(s.strgFileName)
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

func (s *FileStorage) shortIsDeleted(ctx context.Context, short models.ShortURL) (bool, error) {
	s.delMu.Lock()
	defer s.delMu.Unlock()

	fd, err := newFileDecoder(s.delFileName)
	if err != nil {
		return false, err
	}
	defer fd.close()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		deleted := &models.DelURLReq{}
		err := fd.Decode(deleted)
		if err == io.EOF {
			return false, nil
		}
		if err != nil {
			return false, err
		}

		if deleted.Short == short {
			return true, nil
		}
	}
}

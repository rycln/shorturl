package storage

import (
	"context"
	"encoding/json"
	"fmt"
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

func (s *FileStorage) getPairByShort(ctx context.Context, short models.ShortURL) (pair *models.URLPair, err error) {
	s.strgMu.Lock()
	defer s.strgMu.Unlock()

	fd, err := newFileDecoder(s.strgFileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if decCloseErr := fd.close(); decCloseErr != nil {
			err = fmt.Errorf("%v; decoder close failed: %w", err, decCloseErr)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		pair = &models.URLPair{}
		err = fd.Decode(pair)
		if err == io.EOF {
			return nil, errNotExist
		}
		if err != nil {
			return nil, err
		}

		if pair.Short == short {
			return pair, nil
		}
	}
}

func (s *FileStorage) getAllUserPairs(ctx context.Context, uid models.UserID) (userpairs []models.URLPair, err error) {
	s.strgMu.Lock()
	defer s.strgMu.Unlock()

	fd, err := newFileDecoder(s.strgFileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if decCloseErr := fd.close(); decCloseErr != nil {
			err = fmt.Errorf("%v; decoder close failed: %w", err, decCloseErr)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		pair := &models.URLPair{}
		err = fd.Decode(pair)
		if err == io.EOF {
			if userpairs == nil {
				return nil, errNotExist
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

func (s *FileStorage) shortIsDeleted(ctx context.Context, short models.ShortURL) (isDeleted bool, err error) {
	s.delMu.Lock()
	defer s.delMu.Unlock()

	fd, err := newFileDecoder(s.delFileName)
	if err != nil {
		return false, err
	}
	defer func() {
		if decCloseErr := fd.close(); decCloseErr != nil {
			err = fmt.Errorf("%v; decoder close failed: %w", err, decCloseErr)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		deleted := &models.DelURLReq{}
		err = fd.Decode(deleted)
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

package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rycln/shorturl/internal/models"
)

type fileEncoder struct {
	*json.Encoder
	file *os.File
}

func newFileEncoder(fileName string) (*fileEncoder, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &fileEncoder{
		file:    file,
		Encoder: json.NewEncoder(file),
	}, nil
}

func (f *fileEncoder) close() error {
	return f.file.Close()
}

func (s *FileStorage) writeIntoStrgFile(pair *models.URLPair) (err error) {
	s.strgMu.Lock()
	defer s.strgMu.Unlock()

	enc, err := newFileEncoder(s.strgFileName)
	if err != nil {
		return err
	}
	defer func() {
		if encCloseErr := enc.close(); encCloseErr != nil {
			err = fmt.Errorf("%v; encoder close failed: %w", err, encCloseErr)
		}
	}()

	return enc.Encode(pair)
}

func (s *FileStorage) writeIntoDelFile(delReq *models.DelURLReq) (err error) {
	s.delMu.Lock()
	defer s.delMu.Unlock()

	enc, err := newFileEncoder(s.delFileName)
	if err != nil {
		return err
	}
	defer func() {
		if encCloseErr := enc.close(); encCloseErr != nil {
			err = fmt.Errorf("%v; encoder close failed: %w", err, encCloseErr)
		}
	}()

	return enc.Encode(delReq)
}

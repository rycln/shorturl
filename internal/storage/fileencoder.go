package storage

import (
	"encoding/json"
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

func (s *FileStorage) writeIntoFile(pair *models.URLPair) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.enc.Encode(pair)
}

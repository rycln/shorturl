package storage

import (
	"encoding/json"
	"os"
)

type fileEncoder struct {
	file    *os.File
	encoder *json.Encoder
}

func newFileEncoder(fileName string) (*fileEncoder, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &fileEncoder{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (fe *fileEncoder) close() error {
	return fe.file.Close()
}

func (fe *fileEncoder) writeIntoFile(surl *ShortenedURL) error {
	return fe.encoder.Encode(surl)
}

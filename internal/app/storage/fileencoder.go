package storage

import (
	"encoding/json"
	"os"
)

type FileEncoder struct {
	file    *os.File
	encoder *json.Encoder
}

func NewFileEncoder(fileName string) (*FileEncoder, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileEncoder{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (fe *FileEncoder) Close() error {
	return fe.file.Close()
}

func (fe *FileEncoder) writeIntoFile(surl *storedURL) error {
	return fe.encoder.Encode(&surl)
}

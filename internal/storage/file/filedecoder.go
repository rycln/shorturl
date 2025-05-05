package storage

import (
	"encoding/json"
	"os"
)

type fileDecoder struct {
	file    *os.File
	decoder *json.Decoder
}

func newFileDecoder(fileName string) (*fileDecoder, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return &fileDecoder{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (fd *fileDecoder) close() error {
	return fd.file.Close()
}

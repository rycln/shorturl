package storage

import (
	"encoding/json"
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

type storager interface {
	AddURL(string, string) bool
}

func (fd *FileDecoder) RestoreStorage(s storager) error {
	for {
		surl := &StoredURL{}
		err := fd.decoder.Decode(&surl)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		s.AddURL(surl.ShortURL, surl.FullURL)
	}
}

package storage

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type StoredURL struct {
	ID       string `json:"uuid"`
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}

func NewStoredURL(shortURL, fullURL string) *StoredURL {
	surl := &StoredURL{
		ID:       uuid.NewString(),
		ShortURL: shortURL,
		FullURL:  fullURL,
	}
	return surl
}

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

func (fe *FileEncoder) WriteInto(surl *StoredURL) error {
	return fe.encoder.Encode(&surl)
}

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
			return err
		}
		s.AddURL(surl.ShortURL, surl.FullURL)
	}
}

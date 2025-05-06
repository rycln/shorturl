package storage

import (
	"github.com/rycln/shorturl/internal/service"
)

type Storage interface {
	service.ShortenerStorage
	service.BatchShortenerStorage
	service.PingStorage
	service.BatchDeleteStorage
	Close()
}

func Factory(cfg *StorageConfig) Storage {
	switch cfg.strgType {
	default:
		return NewAppMemStorage()
	}
}

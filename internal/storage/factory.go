package storage

import (
	"github.com/rycln/shorturl/internal/service"
	"github.com/rycln/shorturl/internal/storage/appmem"
	"github.com/rycln/shorturl/internal/storage/db"
)

type Storage interface {
	Close()
	service.ShortenerStorage
	service.BatchShortenerStorage
	service.PingStorage
	service.BatchDeleteStorage
}

func Factory(cfg *StorageConfig) Storage {
	switch cfg.strgType {
	case "db":
		return db.NewDatabaseStorage(cfg.databaseDsn)
	default:
		return appmem.NewSimpleStorage()
	}
}

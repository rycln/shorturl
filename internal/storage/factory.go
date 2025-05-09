package storage

import (
	"github.com/rycln/shorturl/internal/service"
)

type Storage interface {
	service.ShortenerStorage
	service.BatchShortenerStorage
	service.PingStorage
	service.BatchDeleterStorage
	Close()
}

func StorageFactory(cfg *StorageConfig) (Storage, error) {
	switch cfg.strgType {
	case "db":
		db, err := NewDB(cfg.databaseDsn)
		if err != nil {
			return nil, err
		}
		return NewDatabaseStorage(db), nil
	case "file":
		return NewFileStorage(cfg.filePath)
	default:
		return NewAppMemStorage(), nil
	}
}

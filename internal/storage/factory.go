package storage

import (
	service "github.com/rycln/shorturl/internal/services"
)

// Storage is an interface returned by a factory to provide concrete storage implementation.
//
// It composes embedded interfaces from service layer and adds an additional Close() method
// for proper resource cleanup. Implementations should ensure all composed service interfaces
// are properly satisfied.
type Storage interface {
	service.ShortenerStorage
	service.BatchShortenerStorage
	service.PingStorage
	service.BatchDeleterStorage
	Close() error
}

// Factory creates a concrete Storage implementation based on the provided type.
//
// Supported storage types:
//   - "db":     database-backed storage (requires configured DB connection)
//   - "file":   persistent file-based storage
//   - default:	 application memoory
//
// The factory handles all initialization logic and returns a ready-to-use Storage
// instance that implements all service interfaces
func Factory(cfg *StorageConfig) (Storage, error) {
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

package storage

// StorageConfig contains all configuration parameters needed to initialize a storage implementation.
type StorageConfig struct {
	strgType    string
	filePath    string
	databaseDsn string
}

type option func(*StorageConfig)

// NewStorageConfig creates a new StorageConfig with provided options.
func NewStorageConfig(opts ...option) *StorageConfig {
	var cfg = &StorageConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// WithStorageType sets the storage implementation type.
func WithStorageType(strg string) option {
	return func(cfg *StorageConfig) {
		cfg.strgType = strg
	}
}

// WithFilePath configures the storage directory for file-based implementation.
func WithFilePath(path string) option {
	return func(cfg *StorageConfig) {
		cfg.filePath = path
	}
}

// WithDatabaseDsn sets the connection string for database storage.
func WithDatabaseDsn(dsn string) option {
	return func(cfg *StorageConfig) {
		cfg.databaseDsn = dsn
	}
}

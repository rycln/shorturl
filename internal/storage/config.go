package storage

type StorageConfig struct {
	strgType    string
	filePath    string
	databaseDsn string
}

type option func(*StorageConfig)

func NewStorageConfig(opts ...option) *StorageConfig {
	var cfg = &StorageConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func WithStorageType(strg string) option {
	return func(cfg *StorageConfig) {
		cfg.strgType = strg
	}
}

func WithFilePath(path string) option {
	return func(cfg *StorageConfig) {
		cfg.filePath = path
	}
}

func WithDatabaseDsn(dsn string) option {
	return func(cfg *StorageConfig) {
		cfg.databaseDsn = dsn
	}
}

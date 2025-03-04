package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

const (
	DefaultServerAddr      = ":8080"
	DefaultBaseAddr        = "http://localhost:8080"
	DefaultStorageFileName = "urls"
)

type Cfg struct {
	ServerAddr      string `env:"SERVER_ADDRESS"`
	ShortBaseAddr   string `env:"BASE_URL"`
	StorageFilePath string `env:"FILE_STORAGE_PATH"`
}

func NewCfg() *Cfg {
	cfg := &Cfg{}

	flag.StringVar(&cfg.ServerAddr, "a", DefaultServerAddr, "Address and port to start the server (environment variable SERVER_ADDRESS has higher priority)")
	flag.StringVar(&cfg.ShortBaseAddr, "b", DefaultBaseAddr, "Base address and port for short URL (environment variable BASE_URL has higher priority)")
	flag.StringVar(&cfg.StorageFilePath, "f", DefaultStorageFileName, "URL storage file path (environment variable FILE_STORAGE_PATH has higher priority)")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func (cfg *Cfg) GetServerAddr() string {
	return cfg.ServerAddr
}

func (cfg *Cfg) GetBaseAddr() string {
	return cfg.ShortBaseAddr
}

func (cfg *Cfg) GetFilePath() string {
	return cfg.StorageFilePath
}

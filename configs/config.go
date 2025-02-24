package config

import (
	"flag"
	"strings"

	"github.com/caarlos0/env/v11"
)

const (
	DefaultServerAddr      = ":8080"
	DefaultBaseAddr        = "http://localhost:8080"
	DefaultFileStoragePath = "."
	storageFileName        = "urls"
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
	flag.StringVar(&cfg.StorageFilePath, "f", DefaultFileStoragePath, "URL storage file path (environment variable FILE_STORAGE_PATH has higher priority)")
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
	if strings.HasSuffix(cfg.StorageFilePath, "/") {
		return cfg.StorageFilePath + storageFileName
	}
	return cfg.StorageFilePath + "/" + storageFileName
}

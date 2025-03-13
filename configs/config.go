package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
)

const (
	defaultServerAddr = ":8080"
	defaultBaseAddr   = "http://localhost:8080"
	defultTimeout     = 2
)

type Cfg struct {
	ServerAddr      string `env:"SERVER_ADDRESS"`
	ShortBaseAddr   string `env:"BASE_URL"`
	StorageFilePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
	Timeout         int    `env:"TIMEOUT_DUR"`
}

func NewCfg() *Cfg {
	cfg := &Cfg{}

	flag.StringVar(&cfg.ServerAddr, "a", defaultServerAddr, "Address and port to start the server (environment variable SERVER_ADDRESS has higher priority)")
	flag.StringVar(&cfg.ShortBaseAddr, "b", defaultBaseAddr, "Base address and port for short URL (environment variable BASE_URL has higher priority)")
	flag.StringVar(&cfg.StorageFilePath, "f", "", "URL storage file path (environment variable FILE_STORAGE_PATH has higher priority)")
	flag.StringVar(&cfg.DatabaseDsn, "d", "", "Database connection address (environment variable DATABASE_DSN has higher priority)")
	flag.IntVar(&cfg.Timeout, "t", defultTimeout, "Timeout duration in seconds (environment variable TIMEOUT_DUR has higher priority)")
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

func (cfg *Cfg) GetDatabaseDsn() string {
	return cfg.DatabaseDsn
}

func (cfg *Cfg) StorageIs() string {
	switch {
	case cfg.DatabaseDsn != "":
		return "db"
	case cfg.StorageFilePath != "":
		return "file"
	default:
		return "app"
	}
}

func (cfg *Cfg) TimeoutDuration() time.Duration {
	return time.Duration(cfg.Timeout) * time.Second
}

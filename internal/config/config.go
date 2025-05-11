package config

import (
	"crypto/rand"
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rycln/shorturl/internal/logger"
)

const (
	defaultServerAddr = ":8080"
	defaultBaseAddr   = "http://localhost:8080"
	defaultTimeout    = time.Duration(2) * time.Minute
	defaultKeyLength  = 32
	defaultLogLevel   = "debug"
)

type Cfg struct {
	ServerAddr      string        `env:"SERVER_ADDRESS"`
	ShortBaseAddr   string        `env:"BASE_URL"`
	StorageFilePath string        `env:"FILE_STORAGE_PATH"`
	DatabaseDsn     string        `env:"DATABASE_DSN"`
	Timeout         time.Duration `env:"TIMEOUT_DUR"`
	Key             string        `env:"JWT_KEY"`
	LogLevel        string        `env:"LOG_LEVEL"`
	StorageType     string        `env:"-"`
}

type ConfigBuilder struct {
	cfg *Cfg
	err error
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		cfg: &Cfg{
			ServerAddr:    defaultServerAddr,
			ShortBaseAddr: defaultBaseAddr,
			Timeout:       defaultTimeout,
			LogLevel:      defaultLogLevel,
		},
		err: nil,
	}
}

func (b *ConfigBuilder) WithFlagParsing() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	flag.StringVar(&b.cfg.ServerAddr, "a", b.cfg.ServerAddr, "Address and port to start the server")
	flag.StringVar(&b.cfg.ShortBaseAddr, "b", b.cfg.ShortBaseAddr, "Base address and port for short URL")
	flag.StringVar(&b.cfg.StorageFilePath, "f", b.cfg.StorageFilePath, "URL storage file path")
	flag.StringVar(&b.cfg.DatabaseDsn, "d", b.cfg.DatabaseDsn, "Database connection address")
	flag.DurationVar(&b.cfg.Timeout, "t", b.cfg.Timeout, "Timeout duration in seconds")
	flag.StringVar(&b.cfg.Key, "k", b.cfg.Key, "Key for jwt autorization")
	flag.StringVar(&b.cfg.LogLevel, "l", b.cfg.LogLevel, "Logger level")
	flag.Parse()

	return b
}

func (b *ConfigBuilder) WithEnvParsing() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	err := env.Parse(b.cfg)
	if err != nil {
		b.cfg = nil
		b.err = err
		return b
	}

	return b
}

func (b *ConfigBuilder) WithDefaultJWTKey() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	if b.cfg.Key == "" {
		key, err := generateKey(defaultKeyLength)
		if err != nil {
			b.cfg = nil
			b.err = err
			return b
		}
		b.cfg.Key = key
		logger.Log.Warn("Default JWT key used!")
	}

	return b
}

func generateKey(n int) (string, error) {
	key := make([]byte, n)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return string(key), nil
}

func (b *ConfigBuilder) WithFilePath(filepath string) *ConfigBuilder {
	if b.err != nil {
		return b
	}

	b.cfg.StorageFilePath = filepath

	return b
}

func (b *ConfigBuilder) WithDatabaseDsn(dsn string) *ConfigBuilder {
	if b.err != nil {
		return b
	}

	b.cfg.DatabaseDsn = dsn

	return b
}

func (b *ConfigBuilder) WithStorageType() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	switch {
	case b.cfg.DatabaseDsn != "":
		b.cfg.StorageType = "db"
	case b.cfg.StorageFilePath != "":
		b.cfg.StorageType = "file"
	default:
		b.cfg.StorageType = "app"
	}

	return b
}

func (b *ConfigBuilder) Build() (*Cfg, error) {
	return b.cfg, b.err
}

// Package config provides centralized application configuration management.
package config

import (
	"crypto/rand"
	"encoding/json"
	"os"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/caarlos0/env/v11"
	"github.com/rycln/shorturl/internal/logger"
)

// Config default values
const (
	defaultServerAddr = ":8080"
	defaultBaseAddr   = "http://localhost:8080"
	defaultTimeout    = time.Duration(2) * time.Minute
	defaultKeyLength  = 32
	defaultLogLevel   = "debug"
)

type CfgFile struct {
	// CfgFileName specifies configuration file name
	CfgFileName string `env:"CONFIG"`
}

// Cfg contains all application configuration parameters.
//
// The structure supports loading from multiple sources:
// - Environment variables (primary)
// - Command-line flags (secondary)
// - Default values (fallback)
//
// Tags specify the corresponding environment variable names.
// StorageType is excluded from env vars as it's derived internally.
type Cfg struct {
	// ServerAddr specifies HTTP server listen address (host:port)
	ServerAddr string `json:"server_address" env:"SERVER_ADDRESS"`

	// ShortBaseAddr is the base URL for shortened links (e.g. "https://sh.rt")
	ShortBaseAddr string `json:"base_url" env:"BASE_URL"`

	// StorageFilePath contains path for file-based storage
	StorageFilePath string `json:"file_storage_path" env:"FILE_STORAGE_PATH"`

	// DatabaseDsn specifies database connection string
	DatabaseDsn string `json:"database_dsn" env:"DATABASE_DSN"`

	// Key contains JWT signing key (min 32 bytes recommended)
	Key string `json:"jwt_key" env:"JWT_KEY"`

	// LogLevel sets logging verbosity (debug|info|warn|error)
	LogLevel string `json:"log_level" env:"LOG_LEVEL"`

	// StorageType is derived from other parameters (memory|file|db)
	StorageType string `json:"-" env:"-"`

	// Timeout defines default network operation timeout
	Timeout time.Duration `json:"timeout_dur" env:"TIMEOUT_DUR"`

	// HTTPS flag
	EnableHTTPS bool `json:"enable_https" env:"ENABLE_HTTPS"`
}

// ConfigBuilder implements builder pattern for Cfg.
type ConfigBuilder struct {
	cfg *Cfg
	err error
}

// NewConfigBuilder creates a new configuration builder with default values.
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

// WithConfigFile load configuration values from specified file.
func (b *ConfigBuilder) WithConfigFile() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	var cfgFile = CfgFile{}
	preFs := flag.NewFlagSet("file-config", flag.ContinueOnError)
	preFs.ParseErrorsWhitelist.UnknownFlags = true

	preFs.StringVarP(&cfgFile.CfgFileName, "config", "c", "", "Path to config file")

	err := preFs.Parse(os.Args[1:])
	if err != nil {
		b.cfg = nil
		b.err = err
		return b
	}

	err = env.Parse(&cfgFile)
	if err != nil {
		b.cfg = nil
		b.err = err
		return b
	}

	if cfgFile.CfgFileName != "" {
		err = getCfgFromFile(cfgFile.CfgFileName, b.cfg)
		if err != nil {
			b.cfg = nil
			b.err = err
			return b
		}
	}

	return b
}

func getCfgFromFile(fname string, cfg *Cfg) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return err
	}

	return nil
}

// WithFlagParsing parses command-line flags into configuration.
func (b *ConfigBuilder) WithFlagParsing() *ConfigBuilder {
	if b.err != nil {
		return b
	}

	flag.StringVarP(&b.cfg.ServerAddr, "a", "a", b.cfg.ServerAddr, "Address and port to start the server")
	flag.StringVarP(&b.cfg.ShortBaseAddr, "b", "b", b.cfg.ShortBaseAddr, "Base address and port for short URL")
	flag.StringVarP(&b.cfg.StorageFilePath, "f", "f", b.cfg.StorageFilePath, "URL storage file path")
	flag.StringVarP(&b.cfg.DatabaseDsn, "d", "d", b.cfg.DatabaseDsn, "Database connection address")
	flag.DurationVarP(&b.cfg.Timeout, "t", "t", b.cfg.Timeout, "Timeout duration in seconds")
	flag.StringVarP(&b.cfg.Key, "k", "k", b.cfg.Key, "Key for jwt autorization")
	flag.StringVarP(&b.cfg.LogLevel, "l", "l", b.cfg.LogLevel, "Logger level")
	flag.BoolVarP(&b.cfg.EnableHTTPS, "s", "s", b.cfg.EnableHTTPS, "Enable HTTPS flag")
	flag.Parse()

	return b
}

// WithEnvParsing loads environment variables into configuration.
//
// Uses struct tags to map variables to fields.
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

// WithDefaultJWTKey sets default jwt key.
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

// WithFilePath sets file storage filepath.
func (b *ConfigBuilder) WithFilePath(filepath string) *ConfigBuilder {
	if b.err != nil {
		return b
	}

	b.cfg.StorageFilePath = filepath

	return b
}

// WithServerAddr sets database dsn string.
func (b *ConfigBuilder) WithDatabaseDsn(dsn string) *ConfigBuilder {
	if b.err != nil {
		return b
	}

	b.cfg.DatabaseDsn = dsn

	return b
}

// Build finalizes configuration and validates values.
//
// Performs storage type auto-detection (prioritizes db > file > memory)
//
// Returns error if any required field is invalid.
func (b *ConfigBuilder) Build() (*Cfg, error) {
	if b.err != nil {
		return nil, b.err
	}

	switch {
	case b.cfg.DatabaseDsn != "":
		b.cfg.StorageType = "db"
	case b.cfg.StorageFilePath != "":
		b.cfg.StorageType = "file"
	default:
		b.cfg.StorageType = "app"
	}

	return b.cfg, b.err
}

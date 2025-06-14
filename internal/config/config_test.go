package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testServerAddr  = ":8081"
	testBaseAddr    = "http://test/"
	testFilePath    = "urls"
	testDatabaseDsn = "test_dsn"
	testTimeout     = time.Duration(3) * time.Minute
	testKey         = "secret_key"
	testLoggerLevel = "info"
	testCfgFileName = "testcfg.json"
)

func TestConfigBuilder_WithEnvParsing(t *testing.T) {
	testCfg := &Cfg{
		ServerAddr:      testServerAddr,
		ShortBaseAddr:   testBaseAddr,
		StorageFilePath: testFilePath,
		DatabaseDsn:     testDatabaseDsn,
		Timeout:         testTimeout,
		Key:             testKey,
		LogLevel:        testLoggerLevel,
		StorageType:     "db",
		EnableHTTPS:     true,
	}

	t.Setenv("SERVER_ADDRESS", testCfg.ServerAddr)
	t.Setenv("BASE_URL", testCfg.ShortBaseAddr)
	t.Setenv("FILE_STORAGE_PATH", testCfg.StorageFilePath)
	t.Setenv("DATABASE_DSN", testCfg.DatabaseDsn)
	t.Setenv("TIMEOUT_DUR", testCfg.Timeout.String())
	t.Setenv("JWT_KEY", testCfg.Key)
	t.Setenv("LOG_LEVEL", testCfg.LogLevel)
	t.Setenv("ENABLE_HTTPS", "true")

	t.Run("valid test", func(t *testing.T) {
		cfg, err := NewConfigBuilder().
			WithEnvParsing().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, testCfg, cfg)
	})
}

func TestConfigBuilder_WithDefaultJWTKey(t *testing.T) {
	t.Run("valid test", func(t *testing.T) {
		cfg, err := NewConfigBuilder().
			WithDefaultJWTKey().
			Build()
		assert.NoError(t, err)
		assert.NotEmpty(t, cfg.Key)
	})
}

func TestConfigBuilder_WithFlagParsing(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	testCfg := &Cfg{
		ServerAddr:      testServerAddr,
		ShortBaseAddr:   testBaseAddr,
		StorageFilePath: testFilePath,
		DatabaseDsn:     testDatabaseDsn,
		Timeout:         testTimeout,
		Key:             testKey,
		LogLevel:        testLoggerLevel,
		StorageType:     "db",
		EnableHTTPS:     true,
	}

	t.Run("valid test", func(t *testing.T) {
		os.Args = []string{
			"./shortener",
			"-a=" + testCfg.ServerAddr,
			"-b=" + testCfg.ShortBaseAddr,
			"-f=" + testCfg.StorageFilePath,
			"-d=" + testCfg.DatabaseDsn,
			"-t=" + testCfg.Timeout.String(),
			"-k=" + testCfg.Key,
			"-l=" + testCfg.LogLevel,
			"-s",
		}

		cfg, err := NewConfigBuilder().
			WithFlagParsing().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, testCfg, cfg)
	})
}

func TestConfigBuilder_WithStorageType(t *testing.T) {
	t.Run("app mem storage type", func(t *testing.T) {
		cfg, err := NewConfigBuilder().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, "app", cfg.StorageType)
	})

	t.Run("file storage type", func(t *testing.T) {
		cfg, err := NewConfigBuilder().
			WithFilePath(testFilePath).
			Build()
		assert.NoError(t, err)
		assert.Equal(t, "file", cfg.StorageType)
	})

	t.Run("db storage type", func(t *testing.T) {
		cfg, err := NewConfigBuilder().
			WithFilePath(testFilePath).
			WithDatabaseDsn(testDatabaseDsn).
			Build()
		assert.NoError(t, err)
		assert.Equal(t, "db", cfg.StorageType)
	})
}

func TestConfigBuilder_WithConfigFile(t *testing.T) {
	testCfg := &Cfg{
		ServerAddr:      testServerAddr,
		ShortBaseAddr:   testBaseAddr,
		StorageFilePath: testFilePath,
		DatabaseDsn:     testDatabaseDsn,
		Timeout:         testTimeout,
		Key:             testKey,
		LogLevel:        testLoggerLevel,
		StorageType:     "db",
		EnableHTTPS:     true,
	}

	file, err := os.Create(testCfgFileName)
	require.NoError(t, err)
	defer func() {
		err = file.Close()
		require.NoError(t, err)
		err = os.Remove(testCfgFileName)
		require.NoError(t, err)
	}()

	enc := json.NewEncoder(file)

	err = enc.Encode(&testCfg)
	require.NoError(t, err)

	t.Run("file name from flag", func(t *testing.T) {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()

		os.Args = []string{
			"./shortener",
			"-a=" + testCfg.ServerAddr,
			"-b=" + testCfg.ShortBaseAddr,
			"-f=" + testCfg.StorageFilePath,
			"-d=" + testCfg.DatabaseDsn,
			"-t=" + testCfg.Timeout.String(),
			"-k=" + testCfg.Key,
			"-l=" + testCfg.LogLevel,
			"-s",
			"-c=" + testCfgFileName,
		}

		cfg, err := NewConfigBuilder().
			WithConfigFile().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, testCfg, cfg)
	})

	t.Run("file name from env", func(t *testing.T) {
		t.Setenv("CONFIG", testCfgFileName)

		cfg, err := NewConfigBuilder().
			WithConfigFile().
			Build()
		assert.NoError(t, err)
		assert.Equal(t, testCfg, cfg)
	})
}

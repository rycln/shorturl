package config

import (
	"flag"
	"os"
)

const (
	defaultServerAddr = ":8080"
	defaultBaseAddr   = "http://localhost:8080"
)

type Cfg struct {
	serverAddr    string
	shortBaseAddr string
}

func NewCfg() *Cfg {
	return &Cfg{}
}

func (cfg *Cfg) Init() {
	flag.StringVar(&cfg.serverAddr, "a", defaultServerAddr, "address and port to run server")
	flag.StringVar(&cfg.shortBaseAddr, "b", defaultBaseAddr, "base address and port for short URL")
	flag.Parse()

	var envAddr string
	envAddr = os.Getenv("SERVER_ADDRESS")
	if envAddr != "" {
		cfg.serverAddr = envAddr
	}
	envAddr = os.Getenv("BASE_URL")
	if envAddr != "" {
		cfg.shortBaseAddr = envAddr
	}
}

func (cfg *Cfg) GetServerAddr() string {
	return cfg.serverAddr
}

func (cfg *Cfg) GetBaseAddr() string {
	return cfg.shortBaseAddr
}

type TestCfg Cfg

func NewTestCfg() *TestCfg {
	return &TestCfg{}
}

func (testCfg *TestCfg) Init() {
	testCfg.serverAddr = defaultServerAddr
	testCfg.shortBaseAddr = defaultBaseAddr
}

func (testCfg *TestCfg) GetServerAddr() string {
	return testCfg.serverAddr
}

func (testCfg *TestCfg) GetBaseAddr() string {
	return testCfg.shortBaseAddr
}

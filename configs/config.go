package config

//"flag"

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

const (
	defaultServerAddr = ":8080"
	defaultBaseAddr   = "http://localhost:8080"
)

type Cfg struct {
	ServerAddr    string `env:"SERVER_ADDRESS"`
	ShortBaseAddr string `env:"BASE_URL"`
}

func NewCfg() *Cfg {
	return &Cfg{}
}

func (cfg *Cfg) Init() {
	flag.StringVar(&cfg.ServerAddr, "a", defaultServerAddr, "address and port to run server")
	flag.StringVar(&cfg.ShortBaseAddr, "b", defaultBaseAddr, "base address and port for short URL")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}
}

func (cfg *Cfg) GetServerAddr() string {
	return cfg.ServerAddr
}

func (cfg *Cfg) GetBaseAddr() string {
	return cfg.ShortBaseAddr
}

type TestCfg Cfg

func NewTestCfg() *TestCfg {
	return &TestCfg{}
}

func (testCfg *TestCfg) Init() {
	testCfg.ServerAddr = defaultServerAddr
	testCfg.ShortBaseAddr = defaultBaseAddr
}

func (testCfg *TestCfg) GetServerAddr() string {
	return testCfg.ServerAddr
}

func (testCfg *TestCfg) GetBaseAddr() string {
	return testCfg.ShortBaseAddr
}

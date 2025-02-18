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
	cfg := &Cfg{}

	flag.StringVar(&cfg.ServerAddr, "a", defaultServerAddr, "address and port to run server")
	flag.StringVar(&cfg.ShortBaseAddr, "b", defaultBaseAddr, "base address and port for short URL")
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

type TestCfg Cfg

func NewTestCfg() *TestCfg {
	testCfg := &TestCfg{}
	testCfg.ServerAddr = defaultServerAddr
	testCfg.ShortBaseAddr = defaultBaseAddr
	return testCfg
}

func (testCfg *TestCfg) GetServerAddr() string {
	return testCfg.ServerAddr
}

func (testCfg *TestCfg) GetBaseAddr() string {
	return testCfg.ShortBaseAddr
}

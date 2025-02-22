package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

const (
	DefaultServerAddr = ":8080"
	DefaultBaseAddr   = "http://localhost:8080"
)

type Cfg struct {
	ServerAddr    string `env:"SERVER_ADDRESS"`
	ShortBaseAddr string `env:"BASE_URL"`
}

func NewCfg() *Cfg {
	cfg := &Cfg{}

	flag.StringVar(&cfg.ServerAddr, "a", DefaultServerAddr, "address and port to run server")
	flag.StringVar(&cfg.ShortBaseAddr, "b", DefaultBaseAddr, "base address and port for short URL")
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

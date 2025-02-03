package config

import (
	"flag"
)

type cfg struct {
	serverAddr    string
	shortBaseAddr string
}

var Config cfg

func Init() {
	flag.StringVar(&Config.serverAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&Config.shortBaseAddr, "b", "http://localhost:8080/", "base address and port for short URL")
	flag.Parse()
}

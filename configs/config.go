package config

import (
	"flag"
)

type cfg struct {
	serverAddr    string
	shortBaseAddr string
}

var values cfg

func Init() {
	flag.StringVar(&values.serverAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&values.shortBaseAddr, "b", "http://localhost:8080/", "base address and port for short URL")
	flag.Parse()
}

func GetServerAddr() string {
	return values.serverAddr
}

func GetShortBaseAddr() string {
	return values.shortBaseAddr
}

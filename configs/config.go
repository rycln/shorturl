package config

import (
	"flag"
	"os"
)

type cfg struct {
	serverAddr    string
	shortBaseAddr string
}

var values cfg

func Init() {
	values.serverAddr = os.Getenv("SERVER_ADDRESS")
	if values.serverAddr == "" {
		flag.StringVar(&values.serverAddr, "a", ":8080", "address and port to run server")
	}

	values.shortBaseAddr = os.Getenv("BASE_URL")
	if values.shortBaseAddr == "" {
		flag.StringVar(&values.shortBaseAddr, "b", "http://localhost:8080", "base address and port for short URL")
	}
	flag.Parse()
}

func GetServerAddr() string {
	return values.serverAddr
}

func GetBaseAddr() string {
	return values.shortBaseAddr
}

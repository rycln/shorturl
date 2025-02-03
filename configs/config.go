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
	flag.StringVar(&values.serverAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&values.shortBaseAddr, "b", "http://localhost:8080", "base address and port for short URL")
	flag.Parse()

	var envAddr string
	envAddr = os.Getenv("SERVER_ADDRESS")
	if envAddr != "" {
		values.serverAddr = envAddr
	}
	envAddr = os.Getenv("BASE_URL")
	if envAddr != "" {
		values.shortBaseAddr = envAddr
	}
}

func GetServerAddr() string {
	return values.serverAddr
}

func GetBaseAddr() string {
	return values.shortBaseAddr
}

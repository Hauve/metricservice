package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Address string
}

func LoadServertConfig() *ServerConfig {
	address := flag.String("a", "localhost:8080", "address")
	flag.Parse()

	addrEnv, ok := os.LookupEnv("ADDRESS")
	if ok {
		*address = addrEnv
	}

	return &ServerConfig{
		Address: *address,
	}
}

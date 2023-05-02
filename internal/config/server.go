package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Address         string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
}

func LoadServerConfig() *ServerConfig {
	address := flag.String("a", "localhost:8080", "address")

	storeInterval := flag.Int("i", 300, "store interval")
	filePath := flag.String("f", "/tmp/metrics-db.json", "file storage path")
	restore := flag.Bool("r", true, "restore")
	flag.Parse()

	storeIntervalEnv, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		temp, err := strconv.Atoi(storeIntervalEnv)
		if err == nil {
			*storeInterval = temp
		}
	}

	filePathEnv, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		*filePath = filePathEnv
	}

	restoreEnv, ok := os.LookupEnv("RESTORE")
	if ok {
		temp, err := strconv.ParseBool(restoreEnv)
		if err == nil {
			*restore = temp
		}
	}

	addrEnv, ok := os.LookupEnv("ADDRESS")
	if ok {
		*address = addrEnv
	}

	return &ServerConfig{
		Address:         *address,
		StoreInterval:   time.Duration(*storeInterval) * time.Second,
		FileStoragePath: *filePath,
		Restore:         *restore,
	}
}

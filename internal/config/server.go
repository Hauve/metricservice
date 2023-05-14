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

	DatabaseDSN string
}

func LoadServerConfig() *ServerConfig {
	address := flag.String("a", "localhost:8080", "address")

	storeInterval := flag.Int("i", 300, "store interval")
	filePath := flag.String("f", "/tmp/metrics-db.json", "file storage path")
	restore := flag.Bool("r", true, "restore")

	dbString := flag.String("d", "", "database dsn")
	flag.Parse()

	storeIntervalEnv, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		temp, err := strconv.Atoi(storeIntervalEnv)
		// Если ошибка, то остаётся значение по умолчанию. По этой причине не обрабатываю её
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
		// Если ошибка, то остаётся значение по умолчанию. По этой причине не обрабатываю её
		if err == nil {
			*restore = temp
		}
	}

	addrEnv, ok := os.LookupEnv("ADDRESS")
	if ok {
		*address = addrEnv
	}

	dbEnv, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		*dbString = dbEnv
	}

	return &ServerConfig{
		Address:         *address,
		StoreInterval:   time.Duration(*storeInterval) * time.Second,
		FileStoragePath: *filePath,
		Restore:         *restore,
		DatabaseDSN:     *dbString,
	}
}

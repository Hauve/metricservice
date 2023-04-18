package main

import (
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/server"
	"github.com/Hauve/metricservice.git/internal/storage"
)

func main() {
	cfg := config.LoadServertConfig()
	storage := storage.NewMemStorage()
	serv := server.New(cfg, storage)
	serv.Run()
}

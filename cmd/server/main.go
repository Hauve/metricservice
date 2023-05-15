package main

import (
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/logger"
	"github.com/Hauve/metricservice.git/internal/server"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
)

func main() {

	cfg := config.LoadServerConfig()
	lg, err := logger.New()
	if err != nil {
		log.Fatalf("Logger creating failed: %s", err)
	}

	r := chi.NewRouter()

	storageDB, storageMem, db := storage.GetStorage(cfg.DatabaseDSN)

	var serv *server.MyServer
	if storageDB == nil {
		serv = server.New(cfg, storageMem, r, lg, db)
	} else if storageMem == nil {
		serv = server.New(cfg, storageDB, r, lg, db)
	}

	serv.Run()
}

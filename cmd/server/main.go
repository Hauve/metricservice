package main

import (
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/server"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.LoadServerConfig()

	r := chi.NewRouter()

	storageDB, storageMem, db := storage.GetStorage(cfg.DatabaseDSN)

	var serv *server.MyServer
	if storageDB == nil {
		serv = server.New(cfg, storageMem, r, db)
	} else if storageMem == nil {
		serv = server.New(cfg, storageDB, r, db)
	}

	serv.Run()
}

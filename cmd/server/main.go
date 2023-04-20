package main

import (
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/server"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadServerConfig()
	st := storage.NewMemStorage()
	r := chi.NewRouter()
	serv := server.New(cfg, st, r)
	serv.Run()
}

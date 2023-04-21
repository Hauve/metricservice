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
	st := storage.NewMemStorage()
	r := chi.NewRouter()
	lg, err := logger.New()
	if err != nil {
		log.Printf("Logger creating failed: %s", err)
	}
	serv := server.New(cfg, st, r, lg)
	serv.Run()
}

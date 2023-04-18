package server

import (
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type MyServer struct {
	cfg     *config.ServerConfig
	storage storage.Storage
}

func New(cfg *config.ServerConfig, storage storage.Storage) *MyServer {
	return &MyServer{
		cfg:     cfg,
		storage: storage,
	}
}

func (s *MyServer) Run() {

	r := chi.NewRouter()

	r.Get("/value/{metricType}/{metricName}", s.GetHandler)
	r.Get("/value/{metricType}/{metricName}/", s.GetHandler)
	r.Get("/", s.GetAllHandler)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.PostHandler)
	r.Post("/update/{metricType}/{metricName}/{metricValue}/", s.PostHandler)

	err := http.ListenAndServe(s.cfg.Address, r)
	if err != nil {
		log.Fatalf("cannot ListenAndServe: %s", err)
	}
}

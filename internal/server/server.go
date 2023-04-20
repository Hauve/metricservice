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
	router  chi.Router
}

func New(cfg *config.ServerConfig, storage storage.Storage, router chi.Router) *MyServer {
	return &MyServer{
		cfg:     cfg,
		storage: storage,
		router:  router,
	}
}

func (s *MyServer) Run() {
	s.registerRoutes()
	err := http.ListenAndServe(s.cfg.Address, s.router)
	if err != nil {
		log.Fatalf("cannot ListenAndServe: %s", err)
	}
}

func (s *MyServer) registerRoutes() {
	s.router.Get("/value/{metricType}/{metricName}", s.GetHandler)
	s.router.Get("/value/{metricType}/{metricName}/", s.GetHandler)
	s.router.Get("/", s.GetAllHandler)
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}", s.PostHandler)
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}/", s.PostHandler)
}

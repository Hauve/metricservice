package server

import (
	"github.com/Hauve/metricservice.git/internal/compression"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/logger"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type MyServer struct {
	cfg     *config.ServerConfig
	storage storage.Storage
	router  chi.Router
	logger  logger.Logger
}

func New(cfg *config.ServerConfig, storage storage.Storage, router chi.Router, log *logger.Logger) *MyServer {
	return &MyServer{
		cfg:     cfg,
		storage: storage,
		router:  router,
		logger:  *log,
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
	s.router.Get("/value/{metricType}/{metricName}", s.logger.WithLogging(s.GetHandler))
	s.router.Get("/value/{metricType}/{metricName}/", s.logger.WithLogging(s.GetHandler))
	s.router.Get("/", s.logger.WithLogging(s.GetAllHandler))
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}", s.logger.WithLogging(compression.WithUnpackingGZIP(s.PostHandler)))
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}/", s.logger.WithLogging(compression.WithUnpackingGZIP(s.PostHandler)))

	s.router.Post("/update", s.logger.WithLogging(compression.WithUnpackingGZIP(s.JSONPostHandler)))
	s.router.Post("/update/", s.logger.WithLogging(compression.WithUnpackingGZIP(s.JSONPostHandler)))

	s.router.Post("/value", s.JSONGetHandler)
	s.router.Post("/value/", s.JSONGetHandler)

}

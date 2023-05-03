package server

import (
	"github.com/Hauve/metricservice.git/internal/compression"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/dumper"
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
	dumper  *dumper.Dumper
}

func New(cfg *config.ServerConfig, storage storage.Storage, router chi.Router, log *logger.Logger, dump *dumper.Dumper) *MyServer {
	return &MyServer{
		cfg:     cfg,
		storage: storage,
		router:  router,
		logger:  *log,
		dumper:  dump,
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
	s.router.Get("/value/{metricType}/{metricName}", s.logger.WithLogging(compression.WithGzip(s.GetHandler)))
	s.router.Get("/value/{metricType}/{metricName}/", s.logger.WithLogging(compression.WithGzip(s.GetHandler)))
	s.router.Get("/", s.logger.WithLogging(compression.WithGzip(s.GetAllHandler)))
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}", s.logger.WithLogging(compression.WithGzip(s.PostHandler)))
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}/", s.logger.WithLogging(compression.WithGzip(s.PostHandler)))

	s.router.Post("/update", s.logger.WithLogging(compression.WithGzip(s.JSONPostHandler)))
	s.router.Post("/update/", s.logger.WithLogging(compression.WithGzip(s.JSONPostHandler)))

	s.router.Post("/value", s.logger.WithLogging(compression.WithGzip(s.JSONGetHandler)))
	s.router.Post("/value/", s.logger.WithLogging(compression.WithGzip(s.JSONGetHandler)))

}

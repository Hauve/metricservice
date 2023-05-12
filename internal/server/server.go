package server

import (
	"github.com/Hauve/metricservice.git/internal/compression"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/logger"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	err := s.restore()
	if err != nil {
		s.logger.Fatalf("cannot restore data from file: %s", err)
	}

	s.registerRoutes()

	go s.runDumper()

	err = http.ListenAndServe(s.cfg.Address, s.router)
	if err != nil {
		log.Fatalf("cannot ListenAndServe: %s", err)
	}
}

func (s *MyServer) registerRoutes() {
	s.router.Use(middleware.StripSlashes)
	s.router.Use(s.logger.WithLogging)
	s.router.Use(compression.WithGzip)

	s.router.Get("/", s.GetAllHandler)
	s.router.Route("/update", func(r chi.Router) {
		r.Use(s.dumpToFileMiddleware)
		r.Post("/update/{metricType}/{metricName}/{metricValue}", s.PostHandler)
		r.Post("/update", s.JSONPostHandler)
	})

	s.router.Get("/value/{metricType}/{metricName}", s.GetHandler)
	s.router.Post("/value", s.JSONGetHandler)

}

package server

import (
	"database/sql"
	"github.com/Hauve/metricservice.git/internal/compression"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/logger"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
)

type MyServer struct {
	cfg     *config.ServerConfig
	storage storage.Storage
	router  chi.Router

	db *sql.DB
}

func New(cfg *config.ServerConfig, storage storage.Storage, router chi.Router, db *sql.DB) *MyServer {
	return &MyServer{
		cfg:     cfg,
		storage: storage,
		router:  router,
		db:      db,
	}
}

func (s *MyServer) Run() {
	if err := s.restore(); err != nil {
		logger.Log.Fatalf("cannot restore data from file: %s", err)
	}
	go s.runDumper()

	s.registerRoutes()

	if err := http.ListenAndServe(s.cfg.Address, s.router); err != nil {
		log.Fatalf("cannot ListenAndServe: %s", err)
	}
}

func (s *MyServer) registerRoutes() {
	s.router.Use(middleware.StripSlashes)
	s.router.Use(logger.Log.WithLogging)
	s.router.Use(compression.WithGzip)

	s.router.Get("/", s.GetAllHandler)
	s.router.Route("/update", func(r chi.Router) {
		r.Use(s.dumpToFileMiddleware)
		r.Post("/{metricType}/{metricName}/{metricValue}", s.PostHandler)
		r.Post("/", s.JSONPostHandler)
	})

	s.router.Get("/value/{metricType}/{metricName}", s.GetHandler)
	s.router.Post("/value", s.JSONGetHandler)

	s.router.Get("/ping", s.Ping)

}

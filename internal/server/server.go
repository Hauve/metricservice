package server

import (
	"context"
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
	"time"
)

type MyServer struct {
	cfg     *config.ServerConfig
	storage storage.Storage
	router  chi.Router
	logger  logger.Logger

	db *sql.DB
}

func New(cfg *config.ServerConfig, storage storage.Storage, router chi.Router, log *logger.Logger, db *sql.DB) *MyServer {
	return &MyServer{
		cfg:     cfg,
		storage: storage,
		router:  router,
		logger:  *log,
		db:      db,
	}
}

func (s *MyServer) Run() {
	if err := s.restore(); err != nil {
		s.logger.Fatalf("cannot restore data from file: %s", err)
	}
	go s.runDumper()

	s.registerRoutes()

	if err := http.ListenAndServe(s.cfg.Address, s.router); err != nil {
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
		r.Post("/{metricType}/{metricName}/{metricValue}", s.PostHandler)
		r.Post("/", s.JSONPostHandler)
	})

	s.router.Get("/value/{metricType}/{metricName}", s.GetHandler)
	s.router.Post("/value", s.JSONGetHandler)

	s.router.Get("/ping", s.Ping)

}

func (s *MyServer) Ping(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if s.db == nil {
		s.logger.Info(s.cfg.DatabaseDSN)
		s.logger.Warnf("Ping connect failed: database is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := s.db.PingContext(ctx); err != nil {
		s.logger.Info(s.cfg.DatabaseDSN)
		s.logger.Warnf("Ping connect failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

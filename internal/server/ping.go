package server

import (
	"context"
	"github.com/Hauve/metricservice.git/internal/logger"
	"net/http"
	"time"
)

func (s *MyServer) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if s.db == nil {
		logger.Log.Info(s.cfg.DatabaseDSN)
		logger.Log.Warnf("Ping connect failed: database is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := s.db.PingContext(ctx); err != nil {
		logger.Log.Info(s.cfg.DatabaseDSN)
		logger.Log.Warnf("Ping connect failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

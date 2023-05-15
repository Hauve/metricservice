package server

import (
	"encoding/json"
	"fmt"
	"github.com/Hauve/metricservice.git/internal/logger"
	"net/http"
	"os"
	"time"
)

func (s *MyServer) dump() {
	if s.cfg.StoreInterval == 0 {
		if err := s.dumpToFile(); err != nil {
			logger.Log.Fatalf("cannot dump metric to file: %s", err)
		}
		return
	}
}

func (s *MyServer) runDumper() {
	if s.cfg.StoreInterval > 0 {
		ticker := time.NewTicker(s.cfg.StoreInterval)
		for {
			<-ticker.C
			if err := s.dumpToFile(); err != nil {
				logger.Log.Fatalf("cannot dump metric to file: %s", err)
			}
		}
	}
}

func (s *MyServer) dumpToFile() (err error) {

	file, err := os.OpenFile(s.cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("cannot open or create  file: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			err = fmt.Errorf("cannot close file: %w", err)
		}
	}()

	dataForWriting := s.storage.GetMetrics()

	if err = json.NewEncoder(file).Encode(&dataForWriting); err != nil {
		return fmt.Errorf("cannot encode json to dump file: %w", err)
	}
	return nil
}

func (s *MyServer) dumpToFileMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		s.dump()
	}
	return http.HandlerFunc(fn)
}

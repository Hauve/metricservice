package server

import (
	"encoding/json"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
	"os"
)

func (s *MyServer) restore() {
	if !s.cfg.Restore {
		return
	}
	file, err := os.Open(s.cfg.FileStoragePath)
	if err != nil {
		log.Printf("cant open dump file for restoring: %s", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	metricsFromFile := jsonmodel.Dump{}
	if err = json.NewDecoder(file).Decode(&metricsFromFile); err != nil {
		log.Printf("cannot decode json from dump: %s", err)
		return
	}

	for _, v := range metricsFromFile {
		switch v.MType {
		case storage.Gauge:
			s.storage.SetGauge(v.ID, *v.Value)
		case storage.Counter:
			s.storage.SetCounter(v.ID, *v.Delta)
		default:
			log.Printf("given undefined metric type from dump: %s", v.MType)
			continue
		}
	}
}

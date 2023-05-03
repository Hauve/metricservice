package server

import (
	"encoding/json"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"log"
	"os"
	"time"
)

func (s *MyServer) dump() {
	ticker := time.NewTicker(s.cfg.StoreInterval)

	file, err := os.OpenFile(s.cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("cant open or create dump file: %s", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	for {
		<-ticker.C

		metricsFromFile := jsonmodel.Dump{}

		gaugeNames := s.storage.GetGaugeKeys()
		for _, name := range gaugeNames {
			value, _ := s.storage.GetGauge(name)
			metric := jsonmodel.Metrics{
				ID:    name,
				MType: "gauge",
				Delta: nil,
				Value: &value,
			}
			metricsFromFile = append(metricsFromFile, metric)
		}

		counterNames := s.storage.GetCounterKeys()
		for _, name := range counterNames {
			value, _ := s.storage.GetCounter(name)

			metric := jsonmodel.Metrics{
				ID:    name,
				MType: "counter",
				Delta: &value,
				Value: nil,
			}
			metricsFromFile = append(metricsFromFile, metric)
		}

		if err = json.NewEncoder(file).Encode(&metricsFromFile); err != nil {
			log.Printf("cannot encode json to dump file: %s", err)
			return
		}
	}
}

package server

import (
	"encoding/json"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
	"os"
	"time"
)

func (s *MyServer) dump() {

	ticker := time.NewTicker(s.cfg.StoreInterval)
	for {
		<-ticker.C
		file, err := os.OpenFile(s.cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			log.Printf("cant open or create dump file: %s", err)
			return
		}

		dataForWriting := getAllDataForDump(s.storage)

		if err = json.NewEncoder(file).Encode(&dataForWriting); err != nil {
			log.Printf("cannot encode json to dump file: %s", err)
			return
		}

		_ = file.Close()
	}
}

func getAllDataForDump(store storage.Storage) jsonmodel.Dump {
	metricsForDumping := jsonmodel.Dump{}

	gaugeNames := store.GetGaugeKeys()
	for _, name := range gaugeNames {
		value, _ := store.GetGauge(name)
		metric := jsonmodel.Metrics{
			ID:    name,
			MType: "gauge",
			Delta: nil,
			Value: &value,
		}
		metricsForDumping = append(metricsForDumping, metric)
	}

	counterNames := store.GetCounterKeys()
	for _, name := range counterNames {
		value, _ := store.GetCounter(name)

		metric := jsonmodel.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
			Value: nil,
		}
		metricsForDumping = append(metricsForDumping, metric)
	}

	return metricsForDumping
}

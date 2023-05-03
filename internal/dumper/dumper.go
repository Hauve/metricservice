package dumper

import (
	"encoding/json"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
	"os"
	"time"
)

type Dumper struct {
	file   *os.File
	config *config.ServerConfig
}

func NewDumper(cfg *config.ServerConfig) *Dumper {
	return &Dumper{
		config: cfg,
	}
}

func (dmp Dumper) Dump(st storage.Storage) {
	dumpTicker := time.NewTicker(time.Duration(dmp.config.StoreInterval) * time.Second)

	for {
		select {
		case <-dumpTicker.C:
			dmp.dump(st)
		}
	}
}

func (dmp Dumper) dump(st storage.Storage) {
	file, err := os.OpenFile(dmp.config.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("cant open or create dump file: %s", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	metricsFromFile := jsonmodel.Dump{}

	gaugeNames := st.GetGaugeKeys()
	for _, name := range gaugeNames {
		value, _ := st.GetGauge(name)
		metric := jsonmodel.Metrics{
			ID:    name,
			MType: "gauge",
			Delta: nil,
			Value: &value,
		}
		metricsFromFile.Data = append(metricsFromFile.Data, metric)
	}

	counterNames := st.GetCounterKeys()
	for _, name := range counterNames {
		value, _ := st.GetCounter(name)

		metric := jsonmodel.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
			Value: nil,
		}
		metricsFromFile.Data = append(metricsFromFile.Data, metric)
	}

	if err := json.NewEncoder(file).Encode(&metricsFromFile); err != nil {
		log.Printf("cannot decode json from dump: %s", err)
		return
	}
}

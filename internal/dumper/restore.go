package dumper

import (
	"encoding/json"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
	"os"
)

func (dmp Dumper) Restore(st storage.Storage) {
	if !dmp.config.Restore {
		return
	}
	file, err := os.OpenFile(dmp.config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("cant open or create dump file: %s", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	metricsFromFile := jsonmodel.Dump{}
	if err := json.NewDecoder(file).Decode(&metricsFromFile); err != nil {
		log.Printf("cannot decode json from dump: %s", err)
		return
	}

	for _, v := range metricsFromFile.Data {
		switch v.MType {
		case storage.Gauge:
			st.SetGauge(v.ID, *v.Value)
		case storage.Counter:
			st.SetCounter(v.ID, *v.Delta)
		default:
			log.Printf("Taken undefined metric type from dump")
			continue
		}
	}
}

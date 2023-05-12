package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"io/fs"
	"os"
)

func (s *MyServer) restore() (err error) {
	if !s.cfg.Restore {
		return
	}
	file, err := os.Open(s.cfg.FileStoragePath)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// Чтобы избежать fatal при несуществующем файле для restore
			s.logger.Warnf("file for restoring is not exsists")
			return nil
		}
		return fmt.Errorf("cant open dump file for restoring: %w", err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			err = fmt.Errorf("cannot close file: %w", err)
		}
	}()

	metricsFromFile := jsonmodel.Dump{}
	if err = json.NewDecoder(file).Decode(&metricsFromFile); err != nil {
		return fmt.Errorf("cannot decode json from dump: %w", err)
	}

	for _, v := range metricsFromFile {
		switch v.MType {
		case jsonmodel.Gauge:
			s.storage.SetGauge(v.ID, *v.Value)
		case jsonmodel.Counter:
			s.storage.SetCounter(v.ID, *v.Delta)
		default:
			s.logger.Warnf("given undefined metric type from dump: %s", v.MType)
			continue
		}
	}
	return nil
}

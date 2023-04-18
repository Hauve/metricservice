package agent

import (
	"fmt"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
)

func (ag *MyAgent) sendMetrics() error {
	log.Println("Sending...")
	gaugeNames := ag.storage.GetGaugeKeys()
	for _, name := range *gaugeNames {
		value, _ := ag.storage.GetGauge(name)
		val := fmt.Sprintf("%f", value)
		if err := ag.metricProcessor.Process(name, val, storage.Gauge); err != nil {
			return fmt.Errorf("cannot process counter metric: %w", err)
		}
	}

	counterNames := ag.storage.GetCounterKeys()
	for _, name := range *counterNames {
		value, _ := ag.storage.GetCounter(name)
		val := fmt.Sprintf("%d", value)
		if err := ag.metricProcessor.Process(name, val, storage.Counter); err != nil {
			return fmt.Errorf("cannot process counter metric: %w", err)
		}
	}
	return nil
}

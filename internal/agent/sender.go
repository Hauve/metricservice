package agent

import (
	"fmt"
	"github.com/Hauve/metricservice.git/internal/storage"
)

func (ag *MyAgent) sendMetrics() error {
	gaugeNames := ag.storage.GetGaugeKeys()
	for _, name := range gaugeNames {
		value, _ := ag.storage.GetGauge(name)
		val := fmt.Sprintf("%f", value)
		if err := ag.sender.Send(name, val, storage.Gauge); err != nil {
			return fmt.Errorf("cannot send counter metric: %w", err)
		}
	}

	counterNames := ag.storage.GetCounterKeys()
	for _, name := range counterNames {
		value, _ := ag.storage.GetCounter(name)
		val := fmt.Sprintf("%d", value)
		if err := ag.sender.Send(name, val, storage.Counter); err != nil {
			return fmt.Errorf("cannot send counter metric: %w", err)
		}
	}
	ag.storage.SetCounter("PollCount", 0)
	return nil
}

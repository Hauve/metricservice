package agent

import (
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/processor"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
	"time"
)

type MyAgent struct {
	cfg             *config.AgentConfig
	storage         storage.Storage
	metricProcessor processor.MetricProcessor
}

func New(
	cfg *config.AgentConfig,
	processor processor.MetricProcessor,
	storage storage.Storage,
) *MyAgent {
	return &MyAgent{
		cfg:             cfg,
		storage:         storage,
		metricProcessor: processor,
	}
}

func (ag *MyAgent) Run() {
	pollTicker := time.NewTicker(ag.cfg.PoolInterval)
	reportTicker := time.NewTicker(ag.cfg.ReportInterval)

	for {
		select {
		case <-pollTicker.C:
			ag.collectMetrics()
		case <-reportTicker.C:
			if err := ag.sendMetrics(); err != nil {
				log.Printf("cannot process metrics: %s", err)
			}
		}
	}
}

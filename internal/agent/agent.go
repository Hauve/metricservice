package agent

import (
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/sender"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
	"time"
)

type MyAgent struct {
	cfg     *config.AgentConfig
	storage storage.Storage
	sender  sender.Sender
}

func New(cfg *config.AgentConfig, storage storage.Storage, sender sender.Sender) *MyAgent {
	return &MyAgent{
		cfg:     cfg,
		storage: storage,
		sender:  sender,
	}
}

func (ag *MyAgent) Run() {
	pollTicker := time.NewTicker(ag.cfg.PoolInterval)
	reportTicker := time.NewTicker(ag.cfg.ReportInterval)

	for {
		select {
		case <-pollTicker.C:
			ag.collectMetrics()
			for _, m := range ag.storage.GetMetrics() {
				log.Printf("%v", *m)
			}
		case <-reportTicker.C:
			if err := ag.sendMetrics(); err != nil {
				log.Printf("cannot process metrics: %s", err)
			}
		}
	}
}

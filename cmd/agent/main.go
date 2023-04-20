package main

import (
	"github.com/Hauve/metricservice.git/internal/agent"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/sender"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
)

func main() {
	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}
	storage := storage.NewMemStorage()
	sender := sender.NewSender(cfg)
	ag := agent.New(cfg, storage, sender)
	ag.Run()
}

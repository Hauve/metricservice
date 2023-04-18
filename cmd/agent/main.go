package main

import (
	"github.com/Hauve/metricservice.git/internal/agent"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/processor"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
)

func main() {
	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}
	p := processor.NewSender(cfg)
	s := storage.NewMemStorage()
	ag := agent.New(cfg, p, s)
	ag.Run()
}

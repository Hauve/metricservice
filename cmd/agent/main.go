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
	st := storage.NewMemStorage()
	snd := sender.NewJSONSender(cfg)
	ag := agent.New(cfg, st, snd)
	ag.Run()
}

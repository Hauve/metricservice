package main

import "github.com/Hauve/metricservice.git/internal/agent"

func main() {
	ag := agent.New()
	ag.Run()
}

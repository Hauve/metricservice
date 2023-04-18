package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	Address        string
	PoolInterval   time.Duration
	ReportInterval time.Duration
}

func LoadAgentConfig() (*AgentConfig, error) {
	address := flag.String("a", "localhost:8080", "address")
	reportInterval := flag.Int("r", 3, "Report interval in seconds")
	pollInterval := flag.Int("p", 2, "Poll interval in seconds")
	flag.Parse()

	addrEnv, ok := os.LookupEnv("ADDRESS")
	if ok {
		*address = addrEnv
	}
	repEnv, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		temp, err := strconv.Atoi(repEnv)
		if err != nil {
			return nil, fmt.Errorf("cannot parse ReportInterval: %w", err)
		} else {
			*reportInterval = temp
		}
	}

	pollEnv, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		temp, err := strconv.Atoi(pollEnv)
		if err != nil {
			return nil, fmt.Errorf("cannot parse PollInterval: %w", err)
		} else {
			*pollInterval = temp
		}
	}

	return &AgentConfig{
		Address:        *address,
		PoolInterval:   time.Duration(*pollInterval) * time.Second,
		ReportInterval: time.Duration(*pollInterval) * time.Second,
	}, nil
}

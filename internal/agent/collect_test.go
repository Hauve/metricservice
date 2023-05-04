package agent

import (
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/sender"
	"github.com/Hauve/metricservice.git/internal/storage"
	"testing"
)

func TestMyAgent_collectRuntimeMetrics(t *testing.T) {
	type fields struct {
		cfg     *config.AgentConfig
		storage storage.Storage
		sender  sender.Sender
	}
	tests := []struct {
		name   string
		fields fields
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := &MyAgent{
				cfg:     tt.fields.cfg,
				storage: tt.fields.storage,
				sender:  tt.fields.sender,
			}
			ag.collectRuntimeMetrics()
		})
	}
}

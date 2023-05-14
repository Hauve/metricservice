package config

import (
	"reflect"
	"testing"
	"time"
)

func TestLoadAgentConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *AgentConfig
		wantErr bool
	}{
		{
			name: "ok: default values",
			want: &AgentConfig{
				Address:        "localhost:8080",
				PoolInterval:   2 * time.Second,
				ReportInterval: 10 * time.Second,
			},
			wantErr: false,
		},
	}
	//Test only for default values
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadAgentConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAgentConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadAgentConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

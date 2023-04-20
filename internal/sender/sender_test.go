package sender

import (
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/storage"
	"net/http"
	"testing"
)

func TestSender_Send(t *testing.T) {
	type args struct {
		name  string
		value string
		mt    storage.MetricType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test 1: send gauge",
			args: args{
				name:  "key1",
				value: "55.6",
				mt:    "gauge",
			},
			wantErr: false,
		}, {
			name: "positive test 2: send counter",
			args: args{
				name:  "key1",
				value: "55",
				mt:    "counter",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Sender{
				cfg: &config.AgentConfig{
					Address: "http://localhost:8080",
				},
				client: &http.Client{},
			}
			if err := m.Send(tt.args.name, tt.args.value, tt.args.mt); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package sender

import (
	"bytes"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/storage"
	"io"
	"net/http"
	"testing"
)

type clientTest struct{}

func (c *clientTest) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		Body: io.NopCloser(bytes.NewBuffer([]byte("0000"))),
	}, nil
}

func TestJSONSender_Send(t *testing.T) {

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
			name: "positive test: gauge",
			args: args{
				name:  "MetricName",
				value: "25.1",
				mt:    "gauge",
			},
			wantErr: false,
		},
		{
			name: "positive test: counter",
			args: args{
				name:  "MetricName",
				value: "25",
				mt:    "counter",
			},
			wantErr: false,
		},
		{
			name: "negative test: float as counter value",
			args: args{
				name:  "MetricName",
				value: "25.2",
				mt:    "counter",
			},
			wantErr: true,
		},
		{
			name: "negative test: string with chars as gauge value",
			args: args{
				name:  "MetricName",
				value: "symbols",
				mt:    "gauge",
			},
			wantErr: true,
		},
		{
			name: "negative test: string with chars as counter value",
			args: args{
				name:  "MetricName",
				value: "symbols",
				mt:    "counter",
			},
			wantErr: true,
		},
		{
			name: "negative test: incorrect metrics type",
			args: args{
				name:  "MetricName",
				value: "25",
				mt:    "strange metrics type",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &JSONSender{
				cfg: &config.AgentConfig{
					Address: "localhost:8080",
				},
				client: &clientTest{},
			}
			if err := m.Send(tt.args.name, tt.args.value, tt.args.mt); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

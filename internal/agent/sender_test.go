package agent

import (
	"fmt"
	"github.com/Hauve/metricservice.git/internal/storage"
	"strconv"
	"testing"
)

type senderForTest struct{}

func (s *senderForTest) Send(_, value string, mt storage.MetricType) error {
	if mt != "gauge" && mt != "counter" {
		return fmt.Errorf("bad metric type")
	}
	if _, err := strconv.ParseFloat(value, 64); mt == "gauge" && err != nil {
		return fmt.Errorf("bad value for gauge type: %w", err)
	}
	if _, err := strconv.ParseInt(value, 10, 64); mt == "counter" && err != nil {
		return fmt.Errorf("bad value for gauge type: %w", err)
	}
	return nil
}

func TestMyAgent_sendMetrics(t *testing.T) {
	type sendArgs struct {
		mtype  string
		mvalue string
	}
	tests := []struct {
		name    string
		wantErr bool
		args    sendArgs
	}{
		{
			name:    "positive test: gauge",
			wantErr: false,
			args: sendArgs{
				mtype:  "gauge",
				mvalue: "2.15",
			},
		},
		{
			name:    "positive test: counter",
			wantErr: false,
			args: sendArgs{
				mtype:  "gauge",
				mvalue: "2.15",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := &MyAgent{
				storage: storage.NewMemStorage(),
				sender:  &senderForTest{},
			}
			if tt.args.mtype == "gauge" {
				temp, err := strconv.ParseFloat(tt.args.mvalue, 64)
				if err != nil {
					t.Errorf("Broken mvalue in test %s: %s for %s type", tt.name, tt.args.mvalue, tt.args.mtype)
				}
				ag.storage.SetGauge("name1", temp)
			} else if tt.args.mtype == "counter" {
				temp, err := strconv.ParseInt(tt.args.mvalue, 10, 64)
				if err != nil {
					t.Errorf("Broken mvalue in test %s: %s for %s type", tt.name, tt.args.mvalue, tt.args.mtype)
				}
				ag.storage.SetCounter("name1", temp)
			}
			if err := ag.sendMetrics(); (err != nil) != tt.wantErr {
				t.Errorf("sendMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package agent

import (
	"github.com/Hauve/metricservice.git/internal/storage"
	"reflect"
	"testing"
)

//func TestMyAgent_Run(t *testing.T) {
//	type fields struct {
//		storage        handlers.Storage
//		address        string
//		reportInterval int
//		pollInterval   int
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ag := &MyAgent{
//				storage:        tt.fields.storage,
//				address:        tt.fields.address,
//				reportInterval: tt.fields.reportInterval,
//				pollInterval:   tt.fields.pollInterval,
//			}
//			ag.Run()
//		})
//	}
//}

//func TestMyAgent_collectMetrics(t *testing.T) {
//	type fields struct {
//		storage        handlers.Storage
//		address        string
//		reportInterval int
//		pollInterval   int
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ag := &MyAgent{
//				storage:        tt.fields.storage,
//				address:        tt.fields.address,
//				reportInterval: tt.fields.reportInterval,
//				pollInterval:   tt.fields.pollInterval,
//			}
//			ag.collectMetrics()
//		})
//	}
//}

//func TestMyAgent_sendMetric(t *testing.T) {
//	type fields struct {
//		storage        handlers.Storage
//		address        string
//		reportInterval int
//		pollInterval   int
//	}
//	type args struct {
//		name  string
//		value string
//		mt    storage.MetricType
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ag := &MyAgent{
//				storage:        tt.fields.storage,
//				address:        tt.fields.address,
//				reportInterval: tt.fields.reportInterval,
//				pollInterval:   tt.fields.pollInterval,
//			}
//			ag.sendMetric(tt.args.name, tt.args.value, tt.args.mt)
//		})
//	}
//}

//func TestMyAgent_sendMetrics(t *testing.T) {
//	type fields struct {
//		storage        handlers.Storage
//		address        string
//		reportInterval int
//		pollInterval   int
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ag := &MyAgent{
//				storage:        tt.fields.storage,
//				address:        tt.fields.address,
//				reportInterval: tt.fields.reportInterval,
//				pollInterval:   tt.fields.pollInterval,
//			}
//			ag.sendMetrics()
//		})
//	}
//}

func TestNew(t *testing.T) {
	tt := struct {
		name string
		want *MyAgent
	}{
		name: "Single test New",
		want: &MyAgent{
			storage:        storage.NewMemStorage(),
			address:        "localhost:8080",
			reportInterval: 10,
			pollInterval:   2,
		},
	}
	t.Run(tt.name, func(t *testing.T) {
		if got := New(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("New() = %v, want %v", got, tt.want)
		}
	})
}

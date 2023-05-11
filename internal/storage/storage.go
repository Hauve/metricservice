package storage

import "github.com/Hauve/metricservice.git/internal/jsonmodel"

type Storage interface {
	SetGauge(string, float64)
	GetGauge(string) (*jsonmodel.Metrics, bool)
	AddCounter(string, int64)
	GetCounter(string) (*jsonmodel.Metrics, bool)
	SetCounter(string, int64)

	GetMetrics() []*jsonmodel.Metrics
}

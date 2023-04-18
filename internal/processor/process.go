package processor

import "github.com/Hauve/metricservice.git/internal/storage"

type MetricProcessor interface {
	Process(name, value string, mt storage.MetricType) error
}

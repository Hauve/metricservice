package sender

import "github.com/Hauve/metricservice.git/internal/storage"

type Sender interface {
	Send(name, value string, mt storage.MetricType) error
}

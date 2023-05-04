package sender

import (
	"github.com/Hauve/metricservice.git/internal/storage"
	"net/http"
)

type Sender interface {
	Send(name, value string, mt storage.MetricType) error
}

// interface for testing
type clientDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

package sender

import (
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"net/http"
)

type Sender interface {
	Send(metrics jsonmodel.Metrics) error
}

// interface for testing
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

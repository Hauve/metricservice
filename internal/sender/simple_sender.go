package sender

import (
	"fmt"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/storage"
	"net/http"
)

type SimpleSender struct {
	cfg    *config.AgentConfig
	client clientDoer
}

func NewSimpleSender(cfg *config.AgentConfig) *SimpleSender {
	return &SimpleSender{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (m *SimpleSender) Send(name, value string, mt storage.MetricType) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", m.cfg.Address, mt, name, value)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("cannot create request object: %w", err)
	}
	req.Header.Add("Content-Length", `0`)
	req.Header.Add("Content-Type", `text/plain`)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot create post request: %w", err)
	}
	return resp.Body.Close()

}

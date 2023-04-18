package processor

import (
	"fmt"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/storage"
	"net/http"
)

type sender struct {
	cfg    *config.AgentConfig
	client *http.Client
}

func NewSender(cfg *config.AgentConfig) *sender {
	return &sender{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (m *sender) Process(name, value string, mt storage.MetricType) error {

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

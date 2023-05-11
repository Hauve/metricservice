package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"net/http"
)

type JSONSender struct {
	cfg    *config.AgentConfig
	client httpClient
}

func NewJSONSender(cfg *config.AgentConfig) *JSONSender {
	return &JSONSender{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (m *JSONSender) Send(mt jsonmodel.Metrics) error {

	url := fmt.Sprintf("http://%s/update/", m.cfg.Address)

	encodedData, err := json.Marshal(mt)
	if err != nil {
		return fmt.Errorf("cannot marshal data: %w", err)
	}

	compressedEncodedData, err := compress(encodedData)
	if err != nil {
		return fmt.Errorf("ERROR: cannot compress data: %w", err)
	}

	buf := bytes.NewBuffer(compressedEncodedData)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return fmt.Errorf("cannot create request object: %w", err)
	}

	req.Header.Add("Content-Type", `application/json; charset=utf-8`)
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot create post request: %w", err)
	}
	return resp.Body.Close()

}

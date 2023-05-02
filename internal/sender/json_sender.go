package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"github.com/Hauve/metricservice.git/internal/storage"
	"net/http"
	"strconv"
)

type JSONSender struct {
	cfg    *config.AgentConfig
	client *http.Client
}

func NewJSONSender(cfg *config.AgentConfig) *JSONSender {
	return &JSONSender{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (m *JSONSender) Send(name, value string, mt storage.MetricType) error {

	url := fmt.Sprintf("http://%s/update/", m.cfg.Address)

	jsonData := jsonmodel.Metrics{
		ID:    name,
		MType: mt,
	}
	switch mt {
	case storage.Gauge:
		var temp float64
		temp, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("ERROR: cannot convert gauge metric from string to float64: %w", err)
		}
		jsonData.Value = &temp
	case storage.Counter:
		var temp int64
		temp, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("ERROR: cannot convert counter metric from string to int64: %w", err)
		}
		jsonData.Delta = &temp
	}
	encodedData, err := json.Marshal(jsonData)
	if err != nil {
		return fmt.Errorf("ERROR: cannot marshal data: %w", err)
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

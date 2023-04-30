package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/json_model"
	"github.com/Hauve/metricservice.git/internal/storage"
	"net/http"
	"strconv"
)

type JsonSender struct {
	cfg    *config.AgentConfig
	client *http.Client
}

func NewJSONSender(cfg *config.AgentConfig) *Sender {
	return &Sender{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (m *JsonSender) Send(name, value string, mt storage.MetricType) error {

	url := fmt.Sprintf("http://%s/update/", m.cfg.Address)

	jsonData := json_model.Metrics{
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
		*jsonData.Value = temp
	case storage.Counter:
		var temp int64
		temp, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("ERROR: cannot convert counter metric from string to int64: %w", err)
		}
		*jsonData.Delta = temp
	}
	encodedData, err := json.Marshal(jsonData)
	if err != nil {
		return fmt.Errorf("ERROR: cannot marshal data: %w", err)
	}
	buf := bytes.NewBuffer(encodedData)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return fmt.Errorf("cannot create request object: %w", err)
	}
	req.Header.Add("Content-Type", `application/json; charset=utf-8`)
	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot create post request: %w", err)
	}
	return resp.Body.Close()

}

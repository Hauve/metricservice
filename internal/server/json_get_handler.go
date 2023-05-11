package server

import (
	"encoding/json"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"io"
	"net/http"
	"strings"
	"time"
)

func (s *MyServer) JSONGetHandler(w http.ResponseWriter, r *http.Request) {
	if header := r.Header.Get("Content-Type"); !strings.Contains(header, "application/json") {
		s.logger.Errorf("ERROR: bad content type for current path")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Date", time.Now().String())

	body := r.Body
	buf, err := io.ReadAll(body)
	if err != nil {
		s.logger.Errorf("ERROR: cannot read from body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = body.Close()
	if err != nil {
		s.logger.Errorf("ERROR: cannot close body of request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := jsonmodel.Metrics{}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		s.logger.Errorf("ERROR: cannot unmarshal json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metricType := data.MType
	metricName := data.ID

	var metric *jsonmodel.Metrics
	isMetricFound := false
	switch metricType {
	case jsonmodel.Gauge:
		metric, isMetricFound = s.storage.GetGauge(metricName)
	case jsonmodel.Counter:
		metric, isMetricFound = s.storage.GetCounter(metricName)
	}
	if !isMetricFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	metric.MType = data.MType
	metric.ID = data.ID

	buf, err = json.Marshal(data)
	if err != nil {
		s.logger.Errorf("ERROR: cannot encode data to json in reply: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(buf)
	if err != nil {
		s.logger.Errorf("ERROR: writing fo body is failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

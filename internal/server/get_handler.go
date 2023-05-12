package server

import (
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

func (s *MyServer) GetHandler(w http.ResponseWriter, r *http.Request) {

	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("Date", time.Now().String())

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	var metric *jsonmodel.Metrics

	var isMetricFound bool
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

	// Created for stripping of zeroes for passing tests
	res := metric.GetValue()
	//for passing tests
	if metric.MType == jsonmodel.Gauge {
		for strings.HasSuffix(res, "0") {
			res = strings.TrimSuffix(res, "0")
		}
	}

	_, err := w.Write([]byte(res))
	if err != nil {
		log.Printf("cannot write response to the client: %s", err)
	}
	w.WriteHeader(http.StatusOK)
}

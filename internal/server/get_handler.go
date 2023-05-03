package server

import (
	"fmt"
	"github.com/Hauve/metricservice.git/internal/storage"
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
	var metricValue string
	var isMetricFound bool
	switch metricType {
	case storage.Gauge:
		var val float64
		val, isMetricFound = s.storage.GetGauge(metricName)
		metricValue = fmt.Sprintf("%f", val)
		for strings.HasSuffix(metricValue, "0") {
			metricValue = strings.TrimSuffix(metricValue, "0")
		}
	case storage.Counter:
		var val int64
		val, isMetricFound = s.storage.GetCounter(metricName)
		metricValue = fmt.Sprintf("%d", val)
	}
	if !isMetricFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := w.Write([]byte(metricValue))
	if err != nil {
		log.Printf("cannot write response to the client: %s", err)
	}
	w.WriteHeader(http.StatusOK)
}

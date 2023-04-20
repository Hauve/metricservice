package server

import (
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *MyServer) PostHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("Date", time.Now().String())

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	switch metricType {
	case storage.Gauge:
		valFloat, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			log.Printf("ERROR: cannot parse gauge metric value: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.storage.SetGauge(metricName, valFloat)

	case storage.Counter:
		valInt, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			log.Printf("ERROR: cannot parse counter metric value: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.storage.AddCounter(metricName, valInt)

	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	log.Println(s.storage.GetCounter("PollCount"))
}

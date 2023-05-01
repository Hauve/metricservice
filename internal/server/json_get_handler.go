package server

import (
	"encoding/json"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"github.com/Hauve/metricservice.git/internal/storage"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func (s *MyServer) JSONGetHandler(w http.ResponseWriter, r *http.Request) {
	if header := r.Header.Get("Content-Type"); !strings.Contains(header, "application/json") {
		log.Printf("ERROR: bad content type for current path")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Date", time.Now().String())
	log.Println("g Header set")

	body, err := r.GetBody()
	if err != nil {
		log.Printf("ERROR: cannot get body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	buf, err := io.ReadAll(body)
	if err != nil {
		log.Printf("ERROR: cannot read from body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := jsonmodel.Metrics{}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		log.Printf("ERROR: cannot unmarshal json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("g Metrics Marshaled")

	metricType := data.MType
	metricName := data.MType

	var isMetricFound bool
	switch metricType {
	case storage.Gauge:
		var val float64
		val, isMetricFound = s.storage.GetGauge(metricName)
		*data.Value = val
	case storage.Counter:
		var val int64
		val, isMetricFound = s.storage.GetCounter(metricName)
		*data.Delta = val
	}
	if !isMetricFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Println("g Data in json set")

	buf, err = json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: cannot encode data to json in reply: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		log.Printf("ERROR: writing fo body is failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

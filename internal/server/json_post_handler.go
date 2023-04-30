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

func (s *MyServer) JSONPostHandler(w http.ResponseWriter, r *http.Request) {
	if header := r.Header.Get("Content-Type"); !strings.Contains(header, "application/json") {
		log.Printf("ERROR: bad content type for current path")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Date", time.Now().String())

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

	metricType := data.MType
	metricName := data.MType

	switch metricType {
	case storage.Gauge:
		valFloat := data.Value
		s.storage.SetGauge(metricName, *valFloat)
		_, err = w.Write(buf)
	case storage.Counter:
		valInt := data.Delta
		s.storage.AddCounter(metricName, *valInt)

		*data.Delta, _ = s.storage.GetCounter(metricName)
		buf, err = json.Marshal(data)
		if err != nil {
			log.Printf("ERROR: cannot encode data to json in reply: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(buf)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if err != nil {
		log.Printf("ERROR: writing fo body is failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

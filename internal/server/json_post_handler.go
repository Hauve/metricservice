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

	body := r.Body
	buf, err := io.ReadAll(body)
	if err != nil {
		log.Printf("ERROR: cannot read from body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = body.Close()
	if err != nil {
		log.Printf("cannot close body of request: %s", err)
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

	if data.Value == nil && data.Delta == nil {
		log.Printf("got incorrect json: value and delta are nil")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := data.MType
	metricName := data.ID

	switch metricType {
	case storage.Gauge:
		if data.Value == nil {
			log.Printf("got incorrect json: type='gauge', but value is nil")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valFloat := data.Value
		s.storage.SetGauge(metricName, *valFloat)
	case storage.Counter:
		if data.Delta == nil {
			log.Printf("got incorrect json: type='counter', but delta is nil")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valInt := data.Delta
		s.storage.AddCounter(metricName, *valInt)
		temp, _ := s.storage.GetCounter(metricName)
		data.Delta = &temp
		buf, err = json.Marshal(data)
		if err != nil {
			log.Printf("ERROR: cannot encode data to json in reply: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	_, err = w.Write(buf)
	if err != nil {
		log.Printf("ERROR: writing fo body is failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

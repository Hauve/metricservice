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
	log.Println("JsonPostHandler")
	if header := r.Header.Get("Content-Type"); !strings.Contains(header, "application/json") {
		log.Printf("ERROR: bad content type for current path")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Date", time.Now().String())
	log.Println("p Header set")

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
	log.Println("p Metrics Marshaled")

	metricType := data.MType
	metricName := data.MType

	switch metricType {
	case storage.Gauge:
		valFloat := data.Value
		s.storage.SetGauge(metricName, *valFloat)
	case storage.Counter:
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
	log.Println("p Data in json set")
	_, err = w.Write(buf)
	if err != nil {
		log.Printf("ERROR: writing fo body is failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

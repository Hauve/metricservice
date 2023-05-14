package server

import (
	"encoding/json"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"io"
	"net/http"
	"strings"
	"time"
)

func (s *MyServer) JSONPostHandler(w http.ResponseWriter, r *http.Request) {
	if header := r.Header.Get("Content-Type"); !strings.Contains(header, "application/json") {
		s.logger.Errorf("ERROR: bad content type for current path")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	header := w.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Date", time.Now().String())

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Errorf("ERROR: cannot read from body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = r.Body.Close()
	if err != nil {
		s.logger.Errorf("cannot close body of request: %s", err)
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

	if !data.IsFullFilled() {
		s.logger.Errorf("got incorrect metrics json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricName := data.ID
	switch data.MType {
	case jsonmodel.Gauge:
		s.storage.SetGauge(metricName, *data.Value)
	case jsonmodel.Counter:
		s.storage.AddCounter(metricName, *data.Delta)
		temp, _ := s.storage.GetCounter(metricName)
		data.Delta = temp.Delta
		buf, err = json.Marshal(data)
		if err != nil {
			s.logger.Errorf("ERROR: cannot encode data to json in reply: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = w.Write(buf)
	if err != nil {
		s.logger.Errorf("ERROR: writing fo body is failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

package handlers

import (
	"github.com/Hauve/metricservice.git/internal/storage"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var MyMemStorage = storage.NewMemStorage()

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("Content-Length", "11")
	header.Set("Date", time.Now().String())

	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/update/")
	pathParts := strings.Split(path, "/")

	if len(pathParts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if pathParts[0] != storage.Gauge &&
		pathParts[0] != storage.Counter {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	metric := strings.ToLower(pathParts[0])
	name := pathParts[1]
	val := pathParts[2]

	switch metric {
	case storage.Gauge:
		val, err := strconv.ParseFloat(val, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		MyMemStorage.SetGauge(name, val)

	case storage.Counter:
		val, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		MyMemStorage.SetCounter(name, val)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

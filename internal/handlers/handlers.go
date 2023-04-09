package handlers

import (
	"github.com/Hauve/metricservice.git/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

var MyMemStorage = storage.NewMemStorage()

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/update/")
	pathParts := strings.Split(path, "/")

	if pathParts[0] != storage.Gauge &&
		pathParts[0] != storage.Counter {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if len(pathParts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metric := pathParts[0]
	name := pathParts[1]
	val := pathParts[2]

	switch metric {
	case storage.Gauge:
		_, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ok := MyMemStorage.Set(name, val, storage.Gauge)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case storage.Counter:
		_, err := strconv.ParseFloat(val, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ok := MyMemStorage.Set(name, val, storage.Gauge)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

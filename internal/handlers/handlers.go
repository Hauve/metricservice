package handlers

import (
	"fmt"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var MyMemStorage = storage.NewMemStorage()

func PostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostHandlerStarted")
	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
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
		valFloat, err := strconv.ParseFloat(val, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		MyMemStorage.SetGauge(name, valFloat)

	case storage.Counter:
		valInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		MyMemStorage.SetCounter(name, valInt)

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetHandlerStarted")
	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("Date", time.Now().String())

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	var resValue string
	switch metricType {
	case storage.Gauge:
		val, ok := MyMemStorage.GetGauge(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resValue = fmt.Sprintf("%f", val)
		for strings.HasSuffix(resValue, "0") {
			resValue, _ = strings.CutSuffix(resValue, "0")
		}
	case storage.Counter:
		val, ok := MyMemStorage.GetCounter(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resValue = fmt.Sprintf("%d", val)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := w.Write([]byte(resValue))
	if err != nil {
		return
	}
}

func GetAllHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("GetAllHandlerStarted")

	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Date", time.Now().String())

	templateHtml := `<!DOCTYPE html>
					<html>
						 <head> 
							<meta charset="UTF-8">
							<title>Metrics</title>
						</head>
						<body>`

	gaugeKeys := MyMemStorage.GetGaugeKeys()
	result := "Gauge [keys] = values: <br>"
	for _, gaugeKey := range *gaugeKeys {
		value, _ := MyMemStorage.GetGauge(gaugeKey)
		localRes := fmt.Sprintf("[%s] = %f<br>", gaugeKey, value)
		result += localRes
	}
	result += "<br><br>Counters [keys] = values:<br>"

	counterKeys := MyMemStorage.GetCounterKeys()
	for _, counterKey := range *counterKeys {
		value, _ := MyMemStorage.GetCounter(counterKey)
		localRes := fmt.Sprintf("[%s] = %d<br>", counterKey, value)
		result += localRes
	}

	templateHtml += result
	templateHtml += `	</body>
					</html>`

	_, err := w.Write([]byte(templateHtml))
	if err != nil {
		return
	}
}

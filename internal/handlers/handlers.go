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

type Storage interface {
	SetGauge(key string, val float64)
	GetGauge(key string) (value float64, ok bool)
	SetCounter(key string, val int64)
	GetCounter(key string) (value int64, ok bool)
	GetGaugeKeys() *[]string
	GetCounterKeys() *[]string
}

type Service struct {
	MyMemStorage Storage
}

func New() *Service {
	return &Service{
		storage.NewMemStorage(),
	}
}

func (s *Service) PostHandler(w http.ResponseWriter, r *http.Request) {
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
		s.MyMemStorage.SetGauge(name, valFloat)

	case storage.Counter:
		valInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.MyMemStorage.SetCounter(name, valInt)

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (s *Service) GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetHandlerStarted")
	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("Date", time.Now().String())

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	var resValue string
	switch metricType {
	case storage.Gauge:
		val, ok := s.MyMemStorage.GetGauge(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resValue = fmt.Sprintf("%f", val)
		for strings.HasSuffix(resValue, "0") {
			resValue = strings.TrimSuffix(resValue, "0")
		}
	case storage.Counter:
		val, ok := s.MyMemStorage.GetCounter(metricName)
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

func (s *Service) GetAllHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("GetAllHandlerStarted")

	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Date", time.Now().String())

	templateHTML := `<!DOCTYPE html>
					<html>
						 <head> 
							<meta charset="UTF-8">
							<title>Metrics</title>
						</head>
						<body>`

	gaugeKeys := s.MyMemStorage.GetGaugeKeys()
	result := "Gauge [keys] = values: <br>"
	for _, gaugeKey := range *gaugeKeys {
		value, _ := s.MyMemStorage.GetGauge(gaugeKey)
		localRes := fmt.Sprintf("[%s] = %f<br>", gaugeKey, value)
		result += localRes
	}
	result += "<br><br>Counters [keys] = values:<br>"

	counterKeys := s.MyMemStorage.GetCounterKeys()
	for _, counterKey := range *counterKeys {
		value, _ := s.MyMemStorage.GetCounter(counterKey)
		localRes := fmt.Sprintf("[%s] = %d<br>", counterKey, value)
		result += localRes
	}

	templateHTML += result
	templateHTML += `	</body>
					</html>`

	_, err := w.Write([]byte(templateHTML))
	if err != nil {
		return
	}
}

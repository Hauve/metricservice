package handlers

import (
	"fmt"
	"github.com/Hauve/metricservice.git/internal/storage"
	"log"
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
	NullCounterValue(key string)
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
	log.Println("PostHandlerStarted")
	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("Date", time.Now().String())

	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestURI := strings.TrimPrefix(r.RequestURI, "/update/")
	pathArgs := strings.Split(requestURI, "/")
	if len(pathArgs) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Println(requestURI)
	metricType := pathArgs[0]
	metricName := pathArgs[1]
	metricValue := pathArgs[2]

	switch metricType {
	case storage.Gauge:
		valFloat, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.MyMemStorage.SetGauge(metricName, valFloat)

	case storage.Counter:
		valInt, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.MyMemStorage.SetCounter(metricName, valInt)

	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
}

func (s *Service) GetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetHandlerStarted")
	header := w.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("Date", time.Now().String())

	if r.Method != http.MethodGet {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestURI := strings.TrimPrefix(r.RequestURI, "/value/")
	pathArgs := strings.Split(requestURI, "/")
	if len(pathArgs) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := pathArgs[0]
	metricName := pathArgs[1]

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
	log.Println("GetAllHandlerStarted")

	header := w.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Date", time.Now().String())

	templateHTML := `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Metrics</title></head><body>`

	gaugeKeys := s.MyMemStorage.GetGaugeKeys()
	result := "Gauge [keys] = values: <br>"
	for _, gaugeKey := range *gaugeKeys {
		value, _ := s.MyMemStorage.GetGauge(gaugeKey)
		val := removeZeroes(fmt.Sprintf("%f", value))
		localRes := fmt.Sprintf("[%s] = %s<br>", gaugeKey, val)
		result += localRes
	}
	result += "<br><br>Counters [keys] = values:<br>"

	counterKeys := s.MyMemStorage.GetCounterKeys()
	for _, counterKey := range *counterKeys {
		value, _ := s.MyMemStorage.GetCounter(counterKey)
		val := removeZeroes(fmt.Sprintf("%d", value))
		localRes := fmt.Sprintf("[%s] = %s<br>", counterKey, val)
		result += localRes
	}

	templateHTML += result
	templateHTML += `	</body>
					</html>`

	_, err := w.Write([]byte(templateHTML))
	if err != nil {
		log.Printf("%e", err)
	}
}

func removeZeroes(a string) string {
	for strings.HasSuffix(a, "0") {
		a = strings.TrimSuffix(a, "0")
	}
	return a
}

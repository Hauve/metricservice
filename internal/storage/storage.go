package storage

import "strings"

type MetricType = string

const (
	Gauge   = MetricType("gauge")
	Counter = MetricType("counter")
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (st *MemStorage) SetGauge(key string, val float64) {
	st.gauge[key] = val
}

func (st *MemStorage) GetGauge(key string) (value float64, ok bool) {
	val, ok := st.gauge[key]
	if !ok {
		return float64(0), false
	}
	return val, true
}

func (st *MemStorage) SetCounter(key string, val int64) {
	st.counter[key] = val + st.counter[key]
}

func (st *MemStorage) GetCounter(key string) (value int64, ok bool) {
	val, ok := st.counter[key]
	if !ok {
		return 0, false
	}
	return val, true
}

func (st *MemStorage) GetGaugeKeys() *[]string {
	keys := make([]string, 0)
	for key := range st.gauge {
		key = strings.TrimSpace(key)
		keys = append(keys, key)
	}
	return &keys
}

func (st *MemStorage) GetCounterKeys() *[]string {
	keys := make([]string, 0)
	for key := range st.counter {
		key = strings.TrimSpace(key)
		keys = append(keys, key)
	}
	return &keys
}

func (st *MemStorage) NullCounterValue(key string) {
	st.counter[key] = 0
}

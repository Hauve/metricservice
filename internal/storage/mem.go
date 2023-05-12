package storage

import (
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
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

func (st *MemStorage) GetGauge(key string) (value *jsonmodel.Metrics, ok bool) {
	val, ok := st.gauge[key]
	if !ok {
		return nil, false
	}
	return &jsonmodel.Metrics{
		ID:    key,
		MType: jsonmodel.Gauge,
		Value: &val,
	}, true
}

func (st *MemStorage) AddCounter(key string, val int64) {
	st.counter[key] = val + st.counter[key]
}

func (st *MemStorage) GetCounter(key string) (value *jsonmodel.Metrics, ok bool) {
	val, ok := st.counter[key]
	if !ok {
		return nil, false
	}
	return &jsonmodel.Metrics{
		ID:    key,
		MType: jsonmodel.Counter,
		Delta: &val,
	}, true
}

func (st *MemStorage) GetMetrics() []*jsonmodel.Metrics {
	res := make([]*jsonmodel.Metrics, 0)

	for key := range st.counter {
		if m, ok := st.GetCounter(key); ok {
			res = append(res, m)
		}
	}

	for key := range st.gauge {
		if m, ok := st.GetCounter(key); ok {
			res = append(res, m)
		}
	}
	return res
}

func (st *MemStorage) SetCounter(key string, val int64) {
	st.counter[key] = val
}

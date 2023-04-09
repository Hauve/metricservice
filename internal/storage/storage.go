package storage

import (
	"fmt"
	"strconv"
)

type MetricType = string

const (
	Gauge   MetricType = MetricType("gauge")
	Counter            = MetricType("counter")
)

type MemStorage struct {
	gauge   map[string]string
	counter map[string]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]string),
		counter: make(map[string]string),
	}
}

func (st *MemStorage) Set(key string, val string, metricType MetricType) (ok bool) {
	switch metricType {
	case Counter:
		beforeStr, ok := st.counter[key]
		if !ok {
			st.counter[key] = val
			return true
		}
		beforeInt, err := strconv.ParseInt(beforeStr, 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		valInt, err := strconv.ParseInt(beforeStr, 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		res := fmt.Sprintf("%d", valInt+beforeInt)
		st.counter[key] = res
	case Gauge:
		st.gauge[key] = val
	}
	return true
}

func (st *MemStorage) Get(key string, metricType MetricType) (value string, ok bool) {
	switch metricType {
	case Counter:
		val, ok := st.counter[key]
		if !ok {
			return "", false
		}
		return val, true
	case Gauge:
		val, ok := st.gauge[key]
		if !ok {
			return "", false
		}
		return val, true
	}
	return "", false
}

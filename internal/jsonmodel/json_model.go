package jsonmodel

import "fmt"

const (
	Gauge   = MetricType("gauge")
	Counter = MetricType("counter")
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricType = string

type Dump []Metrics

func (m *Metrics) GetValue() string {
	if m.MType == Gauge {
		return fmt.Sprintf("%f", *m.Value)
	}
	return fmt.Sprintf("%d", *m.Delta)
}

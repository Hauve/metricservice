package storage

const (
	Gauge   = MetricType("gauge")
	Counter = MetricType("counter")
)

type Storage interface {
	SetGauge(key string, val float64)
	GetGauge(key string) (value float64, ok bool)
	AddCounter(key string, val int64)
	GetCounter(key string) (value int64, ok bool)
	GetGaugeKeys() []string
	GetCounterKeys() []string
}

type MetricType = string

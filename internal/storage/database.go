package storage

import (
	"database/sql"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

func (Database) SetGauge(key string, value float64) {
	panic("implement me")
}

func (Database) GetGauge(key string) (*jsonmodel.Metrics, bool) {
	panic("implement me")
}

func (Database) AddCounter(key string, value int64) {
	panic("implement me")
}

func (Database) GetCounter(key string) (*jsonmodel.Metrics, bool) {
	panic("implement me")
}

func (Database) SetCounter(key string, value int64) {
	panic("implement me")
}

func (Database) GetMetrics() []*jsonmodel.Metrics {
	panic("implement me")
}

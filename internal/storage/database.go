package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"time"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, "CREATE TYPE metrics_type AS ENUM ('gauge', 'counter');")
	if err != nil {
		return nil, fmt.Errorf("cannot to create enum metrics_type in database: %w", err)
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx2,
		"CREATE TABLE if NOT EXISTS metrics ("+
			"\"id\" VARCHAR(250) NOT NULL PRIMARY KEY,"+
			"\"mtype\" metrics_type,"+
			"\"delta\" INTEGER,"+
			"\"value\" DOUBLE PRECISION)")
	if err != nil {
		return nil, fmt.Errorf("cannot to create table in database: %w", err)
	}

	return &Database{
		db: db,
	}, nil
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

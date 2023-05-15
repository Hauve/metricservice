package storage

import (
	"database/sql"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"github.com/Hauve/metricservice.git/internal/logger"
)

type Storage interface {
	SetGauge(string, float64)
	GetGauge(string) (*jsonmodel.Metrics, bool)
	AddCounter(string, int64)
	GetCounter(string) (*jsonmodel.Metrics, bool)
	SetCounter(string, int64)

	GetMetrics() []*jsonmodel.Metrics
}

func GetStorage(databaseDSN string) (*MemStorage, *Database, *sql.DB) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		logger.Log.Fatalf("unable to open sql database: %s", err)
	}

	if databaseDSN == "" {
		return NewMemStorage(), nil, db
	} else {
		dbStorage, err := NewDatabase(db)
		if err != nil {
			logger.Log.Fatalf("cannot create database storage: %s", err)
		}
		return nil, dbStorage, db
	}
}

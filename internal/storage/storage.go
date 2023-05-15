package storage

import (
	"database/sql"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"log"
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
	if databaseDSN == "" {
		return NewMemStorage(), nil, nil
	} else {
		db, err := sql.Open("pgx", databaseDSN)
		if err != nil {
			log.Fatalf("unable to open sql database: %s", err)
		}
		return nil, NewDatabase(db), db
	}
}

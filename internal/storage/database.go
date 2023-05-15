package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"github.com/Hauve/metricservice.git/internal/logger"
	"time"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Если ошибка, значит тип уже существует
	_, _ = db.ExecContext(ctx, "CREATE TYPE metrics_type AS ENUM ('gauge', 'counter');")

	ctx2, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx2,
		"CREATE TABLE if NOT EXISTS metrics ("+
			"\"id\" VARCHAR(250) NOT NULL PRIMARY KEY,"+
			"\"mtype\" metrics_type,"+
			"\"delta\" INTEGER,"+
			"\"val\" DOUBLE PRECISION)")
	if err != nil {
		return nil, fmt.Errorf("cannot to create table in database: %w", err)
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) SetGauge(key string, value float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := d.db.ExecContext(ctx, "INSERT INTO metrics (id, mtype, val)"+
		" VALUES(?,?,?);", key, "gauge", value)
	if err != nil {
		logger.Log.Errorf("saving gauge in database is failed: %s", err)
	}
}

func (d *Database) GetGauge(key string) (*jsonmodel.Metrics, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// делаем обращение к db в рамках полученного контекста
	row := d.db.QueryRowContext(ctx, "SELECT val FROM metrics WHERE id = ? AND mtype = 'gauge';", key)

	out := jsonmodel.Metrics{
		ID:    key,
		MType: jsonmodel.Gauge,
		Delta: nil,
		Value: nil,
	}

	if err := row.Scan(&out.Value); err != nil {
		return nil, false
	}

	return &out, true
}

func (d *Database) AddCounter(key string, value int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// делаем обращение к db в рамках полученного контекста
	row := d.db.QueryRowContext(ctx, "SELECT val FROM metrics WHERE id = ? AND mtype = 'counter';", key)

	var delta int64
	if err := row.Scan(&delta); err != nil {
		logger.Log.Errorf("getting counter from database is failed: %s", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()
	_, err := d.db.ExecContext(ctx2, "UPDATE metrics "+
		"SET delta = ? WHERE id = ? AND mtype = 'counter';", delta+value, key)
	if err != nil {
		logger.Log.Errorf("saving counter in database is failed: %s", err)
	}
}

func (d *Database) GetCounter(key string) (*jsonmodel.Metrics, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// делаем обращение к db в рамках полученного контекста
	row := d.db.QueryRowContext(ctx, "SELECT id, mtype, val, delta FROM metrics WHERE id = ? AND mtype = 'counter';", key)

	out := jsonmodel.Metrics{
		ID:    key,
		MType: jsonmodel.Counter,
		Delta: nil,
		Value: nil,
	}

	if err := row.Scan(&out.Delta); err != nil {
		return nil, false
	}

	return &out, true
}

func (d *Database) SetCounter(key string, value int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := d.db.ExecContext(ctx, "INSERT INTO metrics (id, mtype, delta)"+
		" VALUES (?,?,?);", key, "counter", value)
	if err != nil {
		logger.Log.Errorf("saving counter in database is failed: %s", err)
	}
}

func (d *Database) GetMetrics() []*jsonmodel.Metrics {
	metrics := make([]*jsonmodel.Metrics, 0)

	rows, err := d.db.QueryContext(context.TODO(), "SELECT id, mtype, val, delta from metrics;")
	if err != nil {
		logger.Log.Errorf("getting metrics from database is failed: %s", err)
		return nil
	}

	// обязательно закрываем перед возвратом функции
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Log.Errorf("getting all metrics from database is failed: %s", err)
		}
	}(rows)

	// пробегаем по всем записям
	for rows.Next() {
		var v jsonmodel.Metrics
		err = rows.Scan(&v.ID, &v.MType, &v.Value, &v.Delta)
		if err != nil {
			logger.Log.Errorf("getting metrics from database is failed: %s", err)
		}

		metrics = append(metrics, &v)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		logger.Log.Errorf("getting metrics from database is failed: %s", err)
	}
	return metrics
}

package psql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/markphelps/optional"
)

var _ storage.MetricsStorager = (*MetricsStorage)(nil)

type MetricsStorage struct {
	db *sql.DB
}

func NewMetricsStorage(dbPath string) (*MetricsStorage, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("empty database path")
	}

	db, err := sql.Open("pgx", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbStorage := &MetricsStorage{
		db: db,
	}

	if err := dbStorage.initTable(); err != nil {
		return nil, fmt.Errorf("failed to create table in database: %w", err)
	}

	return dbStorage, nil
}

func (s *MetricsStorage) SaveMetric(ctx context.Context, m model.Metric) error {
	switch m.Type {
	case model.KCounter: {
		_, err := s.db.Exec("INSERT INTO counters (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = counters.value + $2", m.Name, m.MustGetInt())
		return err
	}

	case model.KGauge: {
		_, err := s.db.Exec("INSERT INTO gauges (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE set value = $2", m.Name, m.MustGetFloat())
		return err
	}

	default:
		return fmt.Errorf("failed to store metic %s of type %s into database table: unkonw metric type", m.Name, m.Type)
	}

	//if m.Type == model.KCounter {
	//	incValue := m.MustGetInt()
	//	_, err := s.db.Exec("INSERT INTO counters (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = counters.value + $2", m.Name, incValue)
	//	log.Println(err)
	//	return err
	//}
	//
	//if m.Type == model.KGauge {
	//	log.Println("Saving gauge metric")
	//	_, err := s.db.Exec("INSERT INTO gauges (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE set value = $2", m.Name, m.MustGetFloat())
	//	log.Println(err)
	//	return err
	//}
}

func (s *MetricsStorage) GetMetric(ctx context.Context, name string) (model.Metric, error) {
	if counter, err := s.getCounter(name); err != nil {
		res := model.Metric{
			Name:  name,
			Type:  model.KCounter,
			Delta: optional.NewInt64(counter),
		}
		return res, nil
	}

	if gauge, err := s.getGauge(name); err != nil {
		res := model.Metric{
			Name:  name,
			Type:  model.KGauge,
			Value: optional.NewFloat64(gauge),
		}
		return res, nil
	}

	return model.Metric{}, fmt.Errorf("failed to load metric %s from database: no such metric in counters or gauges table", name)
}

func (s *MetricsStorage) Dump(ctx context.Context) error {
	return fmt.Errorf("database storage Dump: no such method for this type of storage")
}

func (s *MetricsStorage) Ping(ctx context.Context) error {
	if s.db == nil {
		return fmt.Errorf("connection error: database is not connected")
	}
	return s.db.Ping()
}

func (s *MetricsStorage) initTable() error {
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS counters (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value BIGINT NOT NULL)")
	if err != nil {
		return fmt.Errorf("failed to create table for counters metrics: %w", err)
	}

	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS gauges (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value DOUBLE PRECISION NOT NULL)")
	if err != nil {
		return fmt.Errorf("failed to create table for gauges metrics: %w", err)
	}

	return nil
}

func (s *MetricsStorage) getGauge(name string) (float64, error) {
	var gauge float64

	err := s.db.QueryRow("select value from gauges where name = $1", name).Scan(&gauge)
	if err != nil {
		return gauge, fmt.Errorf("queryRow failed: %w", err)
	}
	return gauge, nil
}

func (s *MetricsStorage) getCounter(metricName string) (int64, error) {
	var counter int64

	err := s.db.QueryRow("select value from counters where name = $1", metricName).Scan(&counter)
	if err != nil {
		return counter, fmt.Errorf("queryRow failed: %w", err)
	}
	return counter, nil
}

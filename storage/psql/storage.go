package psql

import (
	"database/sql"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var _ storage.MetricsStorager = (*MetricsStorage)(nil)

type MetricsStorage struct {
	db *sql.DB
}

func NewMetricsStorage(dbPath string) (*MetricsStorage, error) {
	db, err := sql.Open("pgx", dbPath)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}

	dbStorage := &MetricsStorage{
		db: db,
	}

	return dbStorage, nil
}

func (s *MetricsStorage) SaveMetric(m model.Metric) {
}

func (s *MetricsStorage) GetMetric(name string) (model.Metric, bool) {
	var dummy model.Metric
	return dummy, true
}

func (s *MetricsStorage) Dump() {
}

func (s *MetricsStorage) Ping() error {
	if s.db == nil {
		err := fmt.Errorf("connection error: database is not connected")
		return err
	}
	return s.db.Ping()
}
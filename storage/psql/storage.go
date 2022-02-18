package psql

import (
	"database/sql"
	"fmt"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
	"github.com/markphelps/optional"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
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
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db in Ctor: %w", err)
	}
	// TODO! добавить тут ping

	dbStorage := &MetricsStorage{
		db: db,
	}

	// TODO! обработать ошибку
	if err := dbStorage.initTable(); err != nil {
		return nil, fmt.Errorf("failed to create table in Ctor: %w", err)
	}

	return dbStorage, nil
}

func (s *MetricsStorage) SaveMetric(m model.Metric) {
	log.Printf("Saving metric %s to db", m.Name)
	// TODO! в psql есть добавление
	// TODO! погуглить как в SQL сделать increment
	if m.Type == model.KCounter {
		// TODO! if not exist -> insert; else update
		// TODO! или через транзакцию
		log.Println("Saving counter metric")
		valueToStore := m.MustGetInt()
		metricFromDB, found := s.GetMetric(m.Name)
		if found {
			valueToStore += metricFromDB.MustGetInt()
		}

		_, err := s.db.Exec("INSERT INTO counters (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET value = $2", m.Name, valueToStore)
		log.Println(err)
		return
	}

	if m.Type == model.KGauge {
		log.Println("Saving gauge metric")
		_, err := s.db.Exec("INSERT INTO gauges (name, value) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE set value = $2", m.Name, m.MustGetFloat())
		log.Println(err)
		return
	}
}

func (s *MetricsStorage) GetMetric(name string) (model.Metric, bool) {
	log.Println("Loading metrics from db")

	if counter, found := s.getCounter(name); found {
		res := model.Metric{
			Name:  name,
			Type:  model.KCounter,
			Delta: optional.NewInt64(counter),
		}
		return res, true
	}

	if gauge, found := s.getGauge(name); found {
		res := model.Metric{
			Name:  name,
			Type:  model.KGauge,
			Value: optional.NewFloat64(gauge),
		}
		return res, true
	}
	log.Println("Loaded metrics from db")

	var dummy model.Metric
	return dummy, false
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

func (s *MetricsStorage) initTable() error {
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS gauges (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value DOUBLE PRECISION NOT NULL)")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS counters (id serial PRIMARY KEY, name VARCHAR (128) UNIQUE NOT NULL, value BIGINT NOT NULL)")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *MetricsStorage) getGauge(name string) (float64, bool) {
	// TODO! возврат ошибок
	log.Println("Getting gauge value from db")
	var gauge float64

	gRows, err := s.db.Query("SELECT name, value FROM gauges")
	if err != nil {
		log.Print(err)
		return gauge, false
	}
	defer gRows.Close()
	for gRows.Next() {
		if err = gRows.Scan(&name, &gauge); err != nil {
			log.Print(err)
			return gauge, false
		}
	}
	if err = gRows.Err(); err != nil {
		log.Print(err)
		return gauge, false
	}

	return gauge, true
}

func (s *MetricsStorage) getCounter(name string) (int64, bool) {
	log.Println("Getting counter value from db")

	var counter int64

	cRows, err := s.db.Query("SELECT name, value FROM counters")
	if err != nil {
		log.Print(err)
		return counter, false
	}
	defer cRows.Close()
	for cRows.Next() {
		if err = cRows.Scan(&name, &counter); err != nil {
			log.Print(err)
			return counter, false
		}
	}
	if err = cRows.Err(); err != nil {
		log.Print(err)
		return counter, false
	}

	return counter, true
}

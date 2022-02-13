package hybrid

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
	"github.com/Carbohz/go-musthave-devops/storage/filebased"
	"github.com/Carbohz/go-musthave-devops/storage/psql"
	"log"
	"time"
)

var _ storage.MetricsStorager = (*MetricsStorage)(nil)

type Config struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	DBPath        string        `env:"DATABASE_DSN"`
}

type MetricsStorage struct {
	config           Config
	fileBasedStorage *filebased.MetricsStorage
	databaseStorage  *psql.MetricsStorage
}

func NewMetricsStorage(config Config) (*MetricsStorage, error) {
	fbsConfig := filebased.Config{
		StoreInterval: config.StoreInterval,
		StoreFile: config.StoreFile,
		Restore: config.Restore,
	}

	fbs, err := filebased.NewMetricsStorage(fbsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create filebased storage in hybrid storage Ctor: %w", err)
	}

	dbs, err := psql.NewMetricsStorage(config.DBPath)
	if err != nil {
		log.Println(fmt.Errorf("failed to create database storage in hybrid storage Ctor: %w", err))
	}

	storage := &MetricsStorage{
		config:           config,
		fileBasedStorage: fbs,
		databaseStorage:  dbs,
	}

	return storage, nil
}

func (s *MetricsStorage) SaveMetric(m model.Metric) {
	if s.databaseStorage != nil {
		s.databaseStorage.SaveMetric(m)
		s.fileBasedStorage.SaveMetric(m)
	} else {
		s.fileBasedStorage.SaveMetric(m)
	}
}

func (s *MetricsStorage) GetMetric(name string) (model.Metric, bool) {
	//if s.databaseStorage != nil {
	//	return s.databaseStorage.GetMetric(name)
	//} else {
	//	return s.fileBasedStorage.GetMetric(name)
	//}
	return s.fileBasedStorage.GetMetric(name)
}

func (s *MetricsStorage) Dump() {
	//if s.databaseStorage != nil {
	//	s.databaseStorage.Dump()
	//} else {
	//	s.fileBasedStorage.Dump()
	//}
	s.fileBasedStorage.Dump()
}

func (s *MetricsStorage) Ping() error {
	if s.databaseStorage != nil {
		return s.databaseStorage.Ping()
	} else {
		return s.fileBasedStorage.Ping()
	}
}
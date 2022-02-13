package hybrid

import (
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
	"github.com/Carbohz/go-musthave-devops/storage/filebased"
	"github.com/Carbohz/go-musthave-devops/storage/psql"
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
	fbs, _ := filebased.NewMetricsStorage(fbsConfig)
	dbs, _ := psql.NewMetricsStorage(config.DBPath)

	storage := &MetricsStorage{
		config:           config,
		fileBasedStorage: fbs,
		databaseStorage:  dbs,
	}

	return storage, nil
}

func (s *MetricsStorage) SaveMetric(m model.Metric) {
	if s.config.DBPath != "" {
		s.databaseStorage.SaveMetric(m)
	} else {
		s.fileBasedStorage.SaveMetric(m)
	}
}

func (s *MetricsStorage) GetMetric(name string) (model.Metric, bool) {
	if s.config.DBPath != "" {
		return s.databaseStorage.GetMetric(name)
	} else {
		return s.fileBasedStorage.GetMetric(name)
	}
}

func (s *MetricsStorage) Dump() {
	if s.config.DBPath != "" {
		s.databaseStorage.Dump()
	} else {
		s.fileBasedStorage.Dump()
	}
}

func (s *MetricsStorage) Ping() error {
	if s.config.DBPath != "" {
		return s.databaseStorage.Ping()
	} else {
		return s.fileBasedStorage.Ping()
	}
}
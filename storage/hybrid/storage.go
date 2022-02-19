package hybrid

import (
	"fmt"
	configsrv "github.com/Carbohz/go-musthave-devops/config/server"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
	"github.com/Carbohz/go-musthave-devops/storage/filebased"
	"github.com/Carbohz/go-musthave-devops/storage/psql"
	"log"
)

var _ storage.MetricsStorager = (*MetricsStorage)(nil)

type MetricsStorage struct {
	config           configsrv.HybridStorageConfig
	fileBasedStorage *filebased.MetricsStorage
	databaseStorage  *psql.MetricsStorage
}

func NewMetricsStorage(config configsrv.HybridStorageConfig) (*MetricsStorage, error) {
	fbsConfig := filebased.Config{
		StoreInterval: config.StoreInterval,
		StoreFile: config.StoreFile,
		Restore: config.Restore,
	}

	var fbs *filebased.MetricsStorage
	dbs, err := psql.NewMetricsStorage(config.DBPath)
	if err != nil {
		log.Println(fmt.Errorf("failed to create database storage in hybrid storage Ctor: %w", err))
		log.Println("Creating fileBased storage instead")

		fbs, err = filebased.NewMetricsStorage(fbsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create filebased storage in hybrid storage Ctor: %w", err)
		}
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
	} else {
		s.fileBasedStorage.SaveMetric(m)
	}
}

func (s *MetricsStorage) GetMetric(name string) (model.Metric, bool) {
	if s.databaseStorage != nil {
		return s.databaseStorage.GetMetric(name)
	} else {
		return s.fileBasedStorage.GetMetric(name)
	}

	//return s.fileBasedStorage.GetMetric(name)
}

func (s *MetricsStorage) Dump() {
	if s.databaseStorage != nil {
		s.databaseStorage.Dump()
	} else {
		s.fileBasedStorage.Dump()
	}

	//s.fileBasedStorage.Dump()
}

func (s *MetricsStorage) Ping() error {
	if s.databaseStorage != nil {
		return s.databaseStorage.Ping()
	} else {
		return s.fileBasedStorage.Ping()
	}
}
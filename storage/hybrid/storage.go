package hybrid

import (
	"context"
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
	fbsConfig := configsrv.FileBasedStorageConfig{
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

func (s *MetricsStorage) SaveMetric(ctx context.Context, m model.Metric) error {
	if s.databaseStorage != nil {
		return s.databaseStorage.SaveMetric(ctx, m)
	} else {
		return s.fileBasedStorage.SaveMetric(ctx, m)
	}
}

func (s *MetricsStorage) GetMetric(ctx context.Context, name string) (model.Metric, error) {
	if s.databaseStorage != nil {
		return s.databaseStorage.GetMetric(ctx, name)
	} else {
		return s.fileBasedStorage.GetMetric(ctx, name)
	}
}

func (s *MetricsStorage) Dump(ctx context.Context) error {
	if s.databaseStorage != nil {
		return s.databaseStorage.Dump(ctx)
	} else {
		return s.fileBasedStorage.Dump(ctx)
	}
}

func (s *MetricsStorage) Ping(ctx context.Context) error {
	if s.databaseStorage != nil {
		return s.databaseStorage.Ping(ctx)
	} else {
		return s.fileBasedStorage.Ping(ctx)
	}
}
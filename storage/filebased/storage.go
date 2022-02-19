package filebased

import (
	"context"
	"encoding/json"
	"fmt"
	configsrv "github.com/Carbohz/go-musthave-devops/config/server"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
	"github.com/Carbohz/go-musthave-devops/storage/inmemory"
	"log"
	"os"
)

var _ storage.MetricsStorager = (*MetricsStorage)(nil)

type MetricsStorage struct {
	inMemoryStorage *inmemory.MetricsStorage
	config configsrv.FileBasedStorageConfig
}

func NewMetricsStorage(config configsrv.FileBasedStorageConfig) (*MetricsStorage, error) {
	inMemoryStorage := inmemory.NewMetricsStorage()

	storage := &MetricsStorage{
		inMemoryStorage: inMemoryStorage,
		config: config,
	}

	if config.Restore {
		if err := storage.LoadMetrics(); err != nil {
			return nil, fmt.Errorf("failed to restore metrics : %w", err)
		}
	}

	return storage, nil
}

func (s *MetricsStorage) SaveMetric(ctx context.Context, m model.Metric) error {
	if err := s.inMemoryStorage.SaveMetric(ctx, m); err != nil {
		return fmt.Errorf("failed to save metric %v in FileBased storage: %w", m, err)
	}
	return nil
}

func (s *MetricsStorage) GetMetric(ctx context.Context, name string) (model.Metric, error) {
	m, err := s.inMemoryStorage.GetMetric(ctx, name)
	if err != nil {
		return model.Metric{}, fmt.Errorf("failed to get metric %v from FileBased storage: %w", m, err)
	}

	return m, nil
}

func (s *MetricsStorage) LoadMetrics() error {
	flag := os.O_RDONLY

	f, err := os.OpenFile(s.config.StoreFile, flag, 0)
	if err != nil {
		return fmt.Errorf("can't open file for loading metrics: %w", err)
	}
	defer f.Close()

	fInfo, err := os.Stat(s.config.StoreFile)
	if err != nil {
		return err
	}

	fSize := fInfo.Size()
	if fSize == 0 {
		return nil
	}

	var metrics map[string]model.Metric

	if err := json.NewDecoder(f).Decode(&metrics); err != nil {
		return fmt.Errorf("can't decode metrics: %w", err)
	}

	ctx := context.Background()
	for _, m := range metrics {
		if err := s.inMemoryStorage.SaveMetric(ctx, m); err != nil {
			return fmt.Errorf("failed to fill inMemory storage with metric %v : %w", m, err)
		}
	}
	//log.Printf("Metrics successfully loaded from file %s", s.config.StoreFile)
	return nil
}

func (s *MetricsStorage) Dump(ctx context.Context) error {
	//log.Printf("Dumping metrics to file %s", s.config.StoreFile)

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(s.config.StoreFile, flag, 0644)
	if err != nil {
		return fmt.Errorf("can't open file for dumping: %w", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	metrics, err := s.inMemoryStorage.GetAllMetrics()
	if err != nil {
		log.Printf("Nothing to dump: %v", err)
		return nil
	}

	if err := encoder.Encode(metrics); err != nil {
		//log.Fatal("Can't encode server's metrics: ", err)
		return fmt.Errorf("can't encode metrics from inMemory storage: %w", err)
	}

	return nil
}

func (s *MetricsStorage) Ping(ctx context.Context) error {
	return fmt.Errorf("FileBased storage ping: no such method for this type of storage")
}
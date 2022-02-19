package filebased

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	configsrv "github.com/Carbohz/go-musthave-devops/config/server"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
	"github.com/Carbohz/go-musthave-devops/storage/inmemory"
	"log"
	"os"
)

var (
  _ storage.MetricsStorager = (*MetricsStorage)(nil)
  errNoFile = errors.New("open file: the system cannot find the path specified")
  errEmptyFile = errors.New("noting to load: empty file")
)

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
			if errors.Is(err, errEmptyFile) || errors.Is(err, errNoFile) {
				log.Printf("Not fatal error: failed to restore metrics : %v", err)
				return storage, nil
			}

			return nil, fmt.Errorf("failed to restore metrics : %w", err)

			//log.Println(err)
			//return storage, nil
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
		//return fmt.Errorf("can't open file for loading metrics: %w", err)
		return errNoFile
	}
	defer f.Close()

	fInfo, err := os.Stat(s.config.StoreFile)
	if err != nil {
		return err
	}

	fSize := fInfo.Size()
	if fSize == 0 {
		//return fmt.Errorf("noting to load: empty file")
		//return nil
		return errEmptyFile
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
	log.Println("Metrics successfully loaded")
	return nil
}

func (s *MetricsStorage) Dump(ctx context.Context) error {
	//log.Printf("Dumping metrics to file %s", s.config.StoreFile)

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(s.config.StoreFile, flag, 0644)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("can't open file for dumping: %w", err)
	}
	defer f.Close()

	metrics, err := s.inMemoryStorage.GetAllMetrics()
	if err != nil {
		log.Printf("Nothing to dump: %v", err)
		return nil
	}

	encoder := json.NewEncoder(f)
	if err := encoder.Encode(metrics); err != nil {
		//log.Fatal("Can't encode server's metrics: ", err)
		return fmt.Errorf("can't encode metrics from inMemory storage: %w", err)
	}

	log.Printf("Metrics successfully stored")
	return nil
}

func (s *MetricsStorage) Ping(ctx context.Context) error {
	return fmt.Errorf("FileBased storage ping: no such method for this type of storage")
}
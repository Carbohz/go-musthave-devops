package filebased

import (
	"encoding/json"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
	"github.com/Carbohz/go-musthave-devops/storage/inmemory"
	"log"
	"os"
	"time"
)

var _ storage.MetricsStorager = (*MetricsStorage)(nil)

type Config struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

type MetricsStorage struct {
	inMemoryStorage *inmemory.MetricsStorage
	config Config
}

func NewMetricsStorage(config Config) (*MetricsStorage, error) {
	inMemoryStorage, _ := inmemory.NewMetricsStorage()

	storage := &MetricsStorage{
		inMemoryStorage: inMemoryStorage,
		config: config,
	}

	if config.Restore {
		storage.LoadMetrics()
	}

	return storage, nil
}

func (s *MetricsStorage) SaveMetric(m model.Metric) {
	s.inMemoryStorage.SaveMetric(m)
}

func (s *MetricsStorage) GetMetric(name string) (model.Metric, bool) {
	return s.inMemoryStorage.GetMetric(name)
}

// TODO! Добавить error
func (s *MetricsStorage) LoadMetrics() {
	log.Printf("Loading metrics from file %s", s.config.StoreFile)

	flag := os.O_RDONLY

	f, err := os.OpenFile(s.config.StoreFile, flag, 0)
	if err != nil {
		log.Print("Can't open file for loading metrics: ", err)
		return
	}
	defer f.Close()

	var metrics map[string]model.Metric

	if err := json.NewDecoder(f).Decode(&metrics); err != nil {
		// TODO! ошибку вместо fatal
		log.Fatal("Can't decode metrics: ", err)
	}

	for _, m := range metrics {
		s.inMemoryStorage.SaveMetric(m)
	}
	log.Printf("Metrics successfully loaded from file %s", s.config.StoreFile)
}

func (s *MetricsStorage) Dump() {
	log.Printf("Dumping metrics to file %s", s.config.StoreFile)

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(s.config.StoreFile, flag, 0644)
	if err != nil {
		log.Fatal("Can't open file for dumping: ", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	if err := encoder.Encode(s.inMemoryStorage.GetAllMetrics()); err != nil {
		log.Fatal("Can't encode server's metrics: ", err)
	}
}

func (s *MetricsStorage) Ping() error {
	return nil
}
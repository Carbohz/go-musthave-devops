package inmemory

import (
	"github.com/markphelps/optional"
	"log"
	"sync"

	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
)

var _ storage.MetricsStorager = (*MetricsStorage)(nil)

type MetricsStorage struct {
	// сделать embeded (без mu)
	mu sync.RWMutex

	metrics map[string]model.Metric
}

func NewMetricsStorage() (*MetricsStorage, error) {
	st := &MetricsStorage{
		metrics: make(map[string]model.Metric),
	}

	return st, nil
}

func (s *MetricsStorage) SaveMetric(m model.Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// вывернуть ( через NotPresent())
	if m.Delta.Present() {
		if v, found := s.metrics[m.Name]; found {
			newValue := m.MustGetInt()
			oldValue := v.MustGetInt()
			log.Printf("Stored counter value was %v, incoming value is %v, so result is %v", oldValue, newValue, oldValue + newValue)
			s.metrics[m.Name] = model.Metric{Name: m.Name, Type: model.KCounter, Delta: optional.NewInt64(oldValue + newValue)}
			return
		}
	}
	//
	s.metrics[m.Name] = m
}

func (s *MetricsStorage) GetMetric(name string) (model.Metric, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// TODO! если not found -> defaultMetric [ {} ]
	v, found := s.metrics[name]
	return v, found
}

// TODO! плохо, раскрываю детали; Лучше создать новую мапу-копию + блокировка
func (s *MetricsStorage) GetAllMetrics() map[string]model.Metric {
	return s.metrics
}

func (s *MetricsStorage) Dump() {
}

func (s *MetricsStorage) Ping() error {
	return nil
}
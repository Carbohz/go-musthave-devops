package inmemory

import (
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
	"sync"
)

var _ storage.MetricsStorager = (*MetricsStorage)(nil)

type MetricsStorage struct {
	mu sync.RWMutex

	gauges map[string]model.GaugeMetric
	counters map[string]model.CounterMetric
}

func NewMetricsStorage() (*MetricsStorage, error) {
	storage := &MetricsStorage{
		gauges: make(map[string]model.GaugeMetric),
		counters: make(map[string]model.CounterMetric),
	}

	return storage, nil
}

func (s *MetricsStorage) SaveGaugeMetric(m model.GaugeMetric) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.gauges[m.Name] = m
}

func (s *MetricsStorage) SaveCounterMetric(m model.CounterMetric) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counters[m.Name] = m
}

func (s *MetricsStorage) GetGaugeMetric(name string) (float64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, found := s.gauges[name]
	return v.Value, found
}

func (s *MetricsStorage) GetCounterMetric(name string) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, found := s.counters[name]
	return v.Value, found
}

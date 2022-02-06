package inmemory

import (
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
	"log"
	"sync"
)

var _ storage.MetricsStorager = (*MetricsStorage)(nil)

type MetricsStorage struct {
	mu sync.RWMutex

	gauges map[string]model.GaugeMetric
	counters map[string]model.CounterMetric
}

func NewMetricsStorage() (*MetricsStorage, error) {
	log.Println("Created NewMetricsStorage")
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
	log.Println("Saved gauge metric")
}

func (s *MetricsStorage) SaveCounterMetric(m model.CounterMetric) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counters[m.Name] = m
	log.Println("Saved counter metric")
}

func (s *MetricsStorage) LoadGaugeMetric(name string) model.GaugeMetric {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.gauges[name]
}

func (s *MetricsStorage) LoadCounterMetric(name string) model.CounterMetric {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.counters[name]
}

package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/storage"
)

var _ storage.MetricsStorager = (*MetricsStorage)(nil)

type MetricsStorage struct {
	mu sync.RWMutex

	metrics map[string]model.Metric
}

func NewMetricsStorage() *MetricsStorage {
	st := &MetricsStorage{
		metrics: make(map[string]model.Metric),
	}

	return st
}

func (s *MetricsStorage) SaveMetric(ctx context.Context, m model.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m.Delta.Present() {
		newDelta := m.MustGetInt()
		if v, found := s.metrics[m.Name]; found {
			oldDelta := v.MustGetInt()
			s.metrics[m.Name] = model.NewCounterMetric(m.Name, oldDelta + newDelta)
			return nil
		} else {
			s.metrics[m.Name] = m
			return nil
		}
	}

	if m.Value.Present() {
		s.metrics[m.Name] = m
		return nil
	}

	return fmt.Errorf("unknown metric type %s was requested to store into inMemory storage", m.Type)
}

func (s *MetricsStorage) GetMetric(ctx context.Context, name string) (model.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, found := s.metrics[name]
	if !found {
		return model.Metric{}, fmt.Errorf("metric %s not found in inMemory storage", name)
	}

	return v, nil
}

func (s *MetricsStorage) GetAllMetrics() map[string]model.Metric {
	s.mu.Lock()
	defer s.mu.Unlock()

	var metricsCopy map[string]model.Metric

	//if len(s.metrics) == 0 {
	//	return metricsCopy
	//}

	for k, v := range s.metrics {
		metricsCopy[k] = v
	}
	return metricsCopy
}

func (s *MetricsStorage) Dump(ctx context.Context) error {
	return fmt.Errorf("InMemory storage dump: no such method for this type of storage")
}

func (s *MetricsStorage) Ping(ctx context.Context) error {
	return fmt.Errorf("InMemory storage ping: no such method for this type of storage")
}
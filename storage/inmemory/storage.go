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
	mu sync.RWMutex

	//gauges   map[string]float64 //map[string]model.Gauge
	//counters map[string]int64
	metrics map[string]model.Metric
}

func NewMetricsStorage() (*MetricsStorage, error) {
	storage := &MetricsStorage{
		//gauges:   make(map[string]float64),
		//counters: make(map[string]int64),
		metrics: make(map[string]model.Metric),
	}

	return storage, nil
}

//func (s *MetricsStorage) SaveGaugeMetric(m model.GaugeMetric) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	s.gauges[m.Name] = m.Value
//}
//
//func (s *MetricsStorage) SaveCounterMetric(m model.CounterMetric) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	s.counters[m.Name] += m.Value
//}

func (s *MetricsStorage) SaveMetric(m model.Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//if m.Delta.Present() {
	//	// Значит это counter метрика
	//	v, found := s.metrics[m.Name]
	//	if found {
	//		// уже есть в хранилище
	//		newValue, _ := m.Delta.Get()
	//		oldValue, _ := v.Delta.Get()
	//		log.Printf("Stored counter value was %v, incoming value is %v, so result is %v", oldValue, newValue, oldValue + newValue)
	//		s.metrics[m.Name] = model.Metric{Name: m.Name, Type: model.KCounter, Delta: optional.NewInt64(oldValue + newValue)}
	//		//v.Delta.Set(oldValue + newValue)
	//	} else {
	//		// новая метрика
	//		s.metrics[m.Name] = m
	//	}
	//} else {
	//	// Значит это gauge метрика
	//	s.metrics[m.Name] = m
	//}

	if m.Delta.Present() {
		// Counter metric
		v, found := s.metrics[m.Name]
		if found {
			newValue := model.MustGetInt(m)
			oldValue := model.MustGetInt(v)
			log.Printf("Stored counter value was %v, incoming value is %v, so result is %v", oldValue, newValue, oldValue + newValue)
			s.metrics[m.Name] = model.Metric{Name: m.Name, Type: model.KCounter, Delta: optional.NewInt64(oldValue + newValue)}
			return
		}
	}

	s.metrics[m.Name] = m
}

func (s *MetricsStorage) GetMetric(name string) (model.Metric, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, found := s.metrics[name]
	return v, found
}

//func (s *MetricsStorage) GetGaugeMetric(name string) (float64, bool) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	v, found := s.gauges[name]
//	return v, found
//}
//
//func (s *MetricsStorage) GetCounterMetric(name string) (int64, bool) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	v, found := s.counters[name]
//	return v, found
//}

package server

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
)

type internalStorage struct {
	GaugeMetrics   map[string]metrics.GaugeMetric
	CounterMetrics map[string]metrics.CounterMetric
}

// insertion into storage
func (s internalStorage) StoreGaugeMetric(name string, value float64) {
	s.GaugeMetrics[name] = metrics.GaugeMetric{
		Base:  metrics.Base{Name: name, Typename: metrics.Gauge},
		Value: value}
}

func (s internalStorage) StoreCounterMetric(name string, value int64) {
	s.CounterMetrics[name] = metrics.CounterMetric{
		Base:  metrics.Base{Name: name, Typename: metrics.Counter},
		Value: s.CounterMetrics[name].Value + value}
}

// storage lookup
func (s internalStorage) FindGaugeMetric(name string) (float64, error) {
	if value, found := s.GaugeMetrics[name]; found {
		return value.Value, nil
	}
	err := fmt.Errorf("Unknown metric \"%s\" of type \"gauge\"", name)
	return -1.0, err
}

func (s internalStorage) FindCounterMetric(name string) (int64, error) {
	if value, found := s.CounterMetrics[name]; found {
		return value.Value, nil
	}
	err := fmt.Errorf("Unknown metric \"%s\" of type \"counter\"", name)
	return -1, err
}

// load data
func (s internalStorage) LoadGaugeMetrics() map[string]metrics.GaugeMetric {
	return s.GaugeMetrics
}

func (s internalStorage) LoadCounterMetrics() map[string]metrics.CounterMetric {
	return s.CounterMetrics
}

// update data
func (s internalStorage) UpdateGaugeMetrics(name string, value float64) {
	s.GaugeMetrics[name] = metrics.GaugeMetric{
		Base:  metrics.Base{Name: name, Typename: metrics.Gauge},
		Value: value}
}

func (s internalStorage) UpdateCounterMetrics(name string, value int64) {
	s.CounterMetrics[name] = metrics.CounterMetric{
		Base:  metrics.Base{Name: name, Typename: metrics.Gauge},
		Value: s.CounterMetrics[name].Value + value}
}


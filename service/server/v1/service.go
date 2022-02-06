package v1

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/Carbohz/go-musthave-devops/storage"
)

var _ server.Processor = (*Service)(nil)

// Сервер должен собирать и хранить произвольные метрики двух типов
type Service struct {
	storage storage.MetricsStorager
}

func NewService(storage storage.MetricsStorager) (*Service, error) {
	svc := &Service{storage: storage}
	return svc, nil
}

func (s *Service) ProcessGaugeMetric(ctx context.Context, m model.GaugeMetric) error {
	s.storage.SaveGaugeMetric(m)
	return nil
}

func (s *Service) ProcessCounterMetric(ctx context.Context, m model.CounterMetric) error {
	s.storage.SaveCounterMetric(m)
	return nil
}

func (s *Service) GetGaugeMetric(name string) (float64, bool) {
	g := s.storage.LoadGaugeMetric(name)
	return g.Value, true
}

func (s *Service) GetCounterMetric(name string) (int64, bool) {
	c := s.storage.LoadCounterMetric(name)
	return c.Value, true
}
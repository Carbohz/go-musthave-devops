package v1

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/Carbohz/go-musthave-devops/storage"
	"log"
)

var _ server.Processor = (*Service)(nil)

// Сервер должен собирать и хранить произвольные метрики двух типов
type Service struct {
	storage storage.MetricsStorager
}

func NewService(storage storage.MetricsStorager) (*Service, error) {
	log.Println("Created NewService")
	svc := &Service{storage: storage}
	return svc, nil
}

// сохраняет gauge метрики в хранилище
func (s *Service) ProcessGaugeMetric(ctx context.Context, m model.GaugeMetric) error {
	s.storage.SaveGaugeMetric(m)
	return nil
}

// сохраняет counter метрики в хранилище
func (s *Service) ProcessCounterMetric(ctx context.Context, m model.CounterMetric) error {
	log.Println("Called ProcessCounterMetric in Service")
	s.storage.SaveCounterMetric(m)
	return nil
}


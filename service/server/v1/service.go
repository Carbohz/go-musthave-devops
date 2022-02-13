package v1

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/Carbohz/go-musthave-devops/service/server"
	"github.com/Carbohz/go-musthave-devops/storage"
)

var _ server.Processor = (*Service)(nil)

type Service struct {
	storage storage.MetricsStorager
}

func NewService(storage storage.MetricsStorager) (*Service, error) {
	svc := &Service{storage: storage}
	return svc, nil
}

func (s *Service) ProcessMetric(ctx context.Context, m model.Metric) error {
	s.storage.SaveMetric(m)
	return nil
}

func (s *Service) GetMetric(name string) (model.Metric, bool) {
	return s.storage.GetMetric(name)
}

func (s *Service) LoadOnStart() {
	s.storage.LoadOnStart()
}

func (s *Service) Dump() {
	s.storage.Dump()
}
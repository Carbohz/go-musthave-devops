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

func NewService(storage storage.MetricsStorager) *Service {
	svc := &Service{storage: storage}
	return svc
}

func (s *Service) SaveMetric(ctx context.Context, m model.Metric) error {
	return s.storage.SaveMetric(ctx, m)
}

func (s *Service) GetMetric(ctx context.Context, name string) (model.Metric, error) {
	return s.storage.GetMetric(ctx, name)
}

func (s *Service) Dump(ctx context.Context) error {
	return s.storage.Dump(ctx)
}

func (s *Service) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}
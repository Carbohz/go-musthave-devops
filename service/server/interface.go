package server

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/model"
)

type Processor interface {
	SaveMetric(ctx context.Context, m model.Metric) error
	GetMetric(ctx context.Context, name string) (model.Metric, error)
	Dump(ctx context.Context) error
	Ping(ctx context.Context) error
}

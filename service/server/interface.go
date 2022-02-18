package server

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/model"
)

type Processor interface {
	ProcessMetric(ctx context.Context, m model.Metric) error
	GetMetric(name string) (model.Metric, bool)
	Dump()
	Ping() error
}

// TODO! добавить model.Valide
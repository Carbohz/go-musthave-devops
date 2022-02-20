//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/model"
)

type MetricsStorager interface {
	SaveMetric(ctx context.Context, m model.Metric) error
	GetMetric(ctx context.Context, name string) (model.Metric, error)
	Dump(ctx context.Context) error
	Ping(ctx context.Context) error
}
package storage

import "github.com/Carbohz/go-musthave-devops/model"

type MetricsStorager interface {
	SaveMetric(m model.Metric)
	GetMetric(name string) (model.Metric, bool)
}
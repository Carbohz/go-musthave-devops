//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import "github.com/Carbohz/go-musthave-devops/model"

type MetricsStorager interface {
	SaveMetric(m model.Metric)
	GetMetric(name string) (model.Metric, bool)
	Dump()
}
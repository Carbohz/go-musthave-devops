//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import "github.com/Carbohz/go-musthave-devops/model"

// TODO! добавить context; возвращение ошибок
type MetricsStorager interface {
	SaveMetric(m model.Metric) // возвращать ошибку
	GetMetric(name string) (model.Metric, bool) // возвращать ошибку
	Dump()// возвращать ошибку
	Ping() error
}
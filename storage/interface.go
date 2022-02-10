package storage

import "github.com/Carbohz/go-musthave-devops/model"

type MetricsStorager interface {
	//SaveGaugeMetric(m model.GaugeMetric)
	//SaveCounterMetric(m model.CounterMetric)
	SaveMetric(m model.Metric)
	//GetGaugeMetric(name string) (float64, bool)
	//GetCounterMetric(name string) (int64, bool)
	GetMetric(name string) (model.Metric, bool)
}
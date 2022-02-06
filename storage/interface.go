package storage

import "github.com/Carbohz/go-musthave-devops/model"

type MetricsStorager interface {
	SaveGaugeMetric(m model.GaugeMetric)
	SaveCounterMetric(m model.CounterMetric)
	LoadGaugeMetric(name string) model.GaugeMetric
	LoadCounterMetric(name string) model.CounterMetric
}
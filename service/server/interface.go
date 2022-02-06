package server

import (
	"context"
	"github.com/Carbohz/go-musthave-devops/model"
)

// Сервер должен собирать и хранить произвольные метрики двух типов:
type Processor interface {
	// сохраняет метрики в хранилище
	ProcessGaugeMetric(ctx context.Context, m model.GaugeMetric) error
	ProcessCounterMetric(ctx context.Context, m model.CounterMetric) error
}

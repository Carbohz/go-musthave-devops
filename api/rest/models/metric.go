package models

import (
	"fmt"
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/markphelps/optional"
)

// если есть тело запроса (например, JSON), то создаем структуру. Иначе излишне

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metrics) ToModelMetric() model.Metric {
	var modelMetric model.Metric
	modelMetric.Name = m.ID
	modelMetric.Type = m.MType

	if m.Delta != nil {
		modelMetric.Delta = optional.NewInt64(*m.Delta)
	} else {
		modelMetric.Value = optional.NewFloat64(*m.Value)
	}

	return modelMetric
}

func FromModelMetrics(modelMetric model.Metric) Metrics {
	var m Metrics
	m.ID = modelMetric.Name
	m.MType = modelMetric.Type

	if modelMetric.Delta.Present() {
		delta := modelMetric.MustGetInt()
		m.Delta = &delta
		return m
	}

	value := modelMetric.MustGetFloat()
	m.Value = &value
	return m
}

func (m Metrics) Validate() error {
	switch m.MType {
	case model.KGauge:
		if m.Value == nil {
			return fmt.Errorf("invalid Value == nil for MType: %s", m.MType)
		}
	case model.KCounter:
		if m.Delta == nil {
			return fmt.Errorf("invalid Delta == nil for MType: %s", m.MType)
		}
	default:
		return fmt.Errorf("unkown MType: %s", m.MType)
	}
	return nil
}

func (m Metrics) String() string {
	if m.Delta != nil {
		return fmt.Sprintf("[ID: %s, MType: %s, Delta: %v, Value: nil]", m.ID, m.MType, *m.Delta)
	}
	return fmt.Sprintf("[ID: %s, MType: %s, Delta: nil, Value: %v]", m.ID, m.MType, *m.Value)
}
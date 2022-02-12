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

func (m Metrics) ToModelMetric() (model.Metric, error) {
	var modelMetric model.Metric
	modelMetric.Name = m.ID
	modelMetric.Type = m.MType

	if m.Delta != nil {
		modelMetric.Delta = optional.NewInt64(*m.Delta)
		return modelMetric, nil
	}

	if m.Value != nil {
		modelMetric.Value = optional.NewFloat64(*m.Value)
		return modelMetric, nil
	}

	err := fmt.Errorf("serialization to model.Metric failed: missing Delta or Value")
	return modelMetric, err
}

func FromModelMetrics(modelMetric model.Metric) (Metrics, error) {
	var m Metrics
	m.ID = modelMetric.Name
	m.MType = modelMetric.Type

	if modelMetric.Delta.Present() {
		delta := modelMetric.MustGetInt()
		m.Delta = &delta
		return m, nil
	}

	if modelMetric.Value.Present() {
		value := modelMetric.MustGetFloat()
		m.Value = &value
		return m, nil
	}

	err := fmt.Errorf("deserialization from model.Metric failed: missing Delta or Value")
	return m, err
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
	if m.Value != nil {
		return fmt.Sprintf("[ID: %s, MType: %s, Delta: nil, Value: %v]", m.ID, m.MType, *m.Value)
	}
	return ""
}
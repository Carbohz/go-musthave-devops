package models

import (
	"github.com/Carbohz/go-musthave-devops/model"
)

// если есть тело запроса (например, JSON), то создаем структуру. Иначе излишне

// представления для api
// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>;

type GaugeMetricRequest struct {
	MType string
	Name  string
	Value float64
}

type GaugeMetricResponse struct {
	//
}

type CounterMetricRequest struct {
	MType string
	Name  string
	Value int64
}

// converter from GaugeMetricRequest to models.GaugeMetric

func (m *GaugeMetricRequest) ToModelGaugeMetric() model.GaugeMetric {
	return model.GaugeMetric{Common: model.Common{Name: m.Name, Typename: m.MType}, Value: m.Value}
}

func (m *CounterMetricRequest) ToModelCounterMetric() model.CounterMetric {
	return model.CounterMetric{Common: model.Common{Name: m.Name, Typename: m.MType}, Value: m.Value}
}

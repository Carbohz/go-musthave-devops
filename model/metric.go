package model

import "strconv"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type (
	Common struct {
		Name string
		Typename string
	}

	GaugeMetric struct {
		Common
		Value float64
	}

	CounterMetric struct {
		Common
		Value int64
	}
)

func (m GaugeMetric) String() string {
	return strconv.FormatFloat(m.Value, 'f', -1, 64)
}

func (m CounterMetric) String() string {
	return strconv.FormatInt(m.Value, 10)
}
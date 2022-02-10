package model

import (
	"strconv"

	"github.com/markphelps/optional"
)


const (
	KGauge   = "gauge"
	KCounter = "counter"
)

type (
	Metric struct {
		Name string
		Type string
		Delta optional.Int64
		Value optional.Float64
	}
)

func (m Metric) String() string {
	delta, err := m.Delta.Get()
	if err == nil {
		return strconv.FormatInt(delta, 10)
	}

	value, _ := m.Value.Get()
	return strconv.FormatFloat(value, 'f', -1, 64)

	//if m.Delta != nil {
	//	return strconv.FormatInt(*m.Delta, 10)
	//}
	//
	//return strconv.FormatFloat(*m.Value, 'f', -1, 64)
}
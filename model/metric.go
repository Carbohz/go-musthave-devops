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

//func MustGetInt(m Metric) int64 {
//	value, err := m.Delta.Get()
//	if err != nil {
//		panic("value not present")
//	}
//	return value
//}

func (m Metric) MustGetInt() int64 {
	value, err := m.Delta.Get()
	if err != nil {
		panic("value not present")
	}
	return value
}

//func MustGetFloat(m Metric) float64 {
//	value, err := m.Value.Get()
//	if err != nil {
//		panic("value not present")
//	}
//	return value
//}

func (m Metric) MustGetFloat() float64 {
	value, err := m.Value.Get()
	if err != nil {
		panic("value not present")
	}
	return value
}

func (m Metric) String() string {
	if m.Delta.Present() {
		delta := m.MustGetInt()
		return strconv.FormatInt(delta, 10)
	}

	value := m.MustGetFloat()
	return strconv.FormatFloat(value, 'f', -1, 64)
}
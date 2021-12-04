package metrics

import (
	"math/rand"
	"runtime"
	"time"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Base struct {
	Name string
	Typename string
}

type GaugeMetric struct {
	Base
	Value float64
}

type CounterMetric struct {
	Base
	Value int64
}

var PollCount = CounterMetric{Base{Name:"PollCount", Typename: Counter}, 0}

func GetRuntimeMetrics() []GaugeMetric {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m := []GaugeMetric{
		{Base{"Alloc",Gauge}, float64(rtm.Alloc)},
		{Base{"BuckHashSys",Gauge}, float64(rtm.BuckHashSys)},
		{Base{"Frees", Gauge},float64(rtm.Frees)},
		{Base{"GCCPUFraction",Gauge}, rtm.GCCPUFraction},
		{Base{"GCSys",Gauge}, float64(rtm.GCSys)},
		{Base{"HeapAlloc",Gauge}, float64(rtm.HeapAlloc)},
		{Base{"HeapIdle",Gauge}, float64(rtm.HeapIdle)},
		{Base{"HeapInuse", Gauge},float64(rtm.HeapInuse)},
		{Base{"HeapObjects",Gauge}, float64(rtm.HeapObjects)},
		{Base{"HeapReleased",Gauge}, float64(rtm.HeapReleased)},
		{Base{"HeapSys",Gauge}, float64(rtm.HeapSys)},
		{Base{"LastGC", Gauge},float64(rtm.LastGC)},
		{Base{"Lookups",Gauge}, float64(rtm.Lookups)},
		{Base{"MCacheInuse", Gauge},float64(rtm.MCacheInuse)},
		{Base{"MCacheSys", Gauge},float64(rtm.MCacheSys)},
		{Base{"MSpanInuse",Gauge},float64(rtm.MSpanInuse)},
		{Base{"MSpanSys",Gauge}, float64(rtm.MSpanSys)},
		{Base{"Mallocs", Gauge},float64(rtm.Mallocs)},
		{Base{"NextGC", Gauge},float64(rtm.NextGC)},
		{Base{"NumForcedGC",Gauge}, float64(rtm.NumForcedGC)},
		{Base{"NumGC",Gauge},float64(rtm.NumGC)},
		{Base{"OtherSys",Gauge}, float64(rtm.OtherSys)},
		{Base{"PauseTotalNs", Gauge},float64(rtm.PauseTotalNs)},
		{Base{"StackInuse",Gauge}, float64(rtm.StackInuse)},
		{Base{"StackSys", Gauge},float64(rtm.StackSys)},
		{Base{"Sys", Gauge},float64(rtm.Sys)},
	}
	return m
}

func GetRandomValueMetric() GaugeMetric {
	return GaugeMetric{Base{"RandomValue", Gauge}, rand.Float64()}
}

func GetPollCountMetric() CounterMetric {
	rand.Seed(time.Now().UnixNano())
	return PollCount
}

func IncrementPollCountMetric() {
	PollCount.Value++
}

func ResetPollCountMetric() {
	PollCount.Value = 0
}
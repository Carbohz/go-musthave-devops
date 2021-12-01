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

type Metric struct {
	Name     string
	Typename string
	Value    float64
}

var PollCount int64 = 0

func GetRuntimeMetrics() []Metric {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	m := []Metric{
		{"Alloc",Gauge, float64(rtm.Alloc)},
		{"BuckHashSys",Gauge, float64(rtm.BuckHashSys)},
		{"Frees", Gauge,float64(rtm.Frees)},
		{"GCCPUFraction",Gauge, rtm.GCCPUFraction},
		{"GCSys",Gauge, float64(rtm.GCSys)},
		{"HeapAlloc",Gauge, float64(rtm.HeapAlloc)},
		{"HeapIdle",Gauge, float64(rtm.HeapIdle)},
		{"HeapInuse", Gauge,float64(rtm.HeapInuse)},
		{"HeapObjects",Gauge, float64(rtm.HeapObjects)},
		{"HeapReleased",Gauge, float64(rtm.HeapReleased)},
		{"HeapSys",Gauge, float64(rtm.HeapSys)},
		{"LastGC", Gauge,float64(rtm.LastGC)},
		{"Lookups",Gauge, float64(rtm.Lookups)},
		{"MCacheInuse", Gauge,float64(rtm.MCacheInuse)},
		{"MCacheSys", Gauge,float64(rtm.MCacheSys)},
		{"MSpanInuse",Gauge,float64(rtm.MSpanInuse)},
		{"MSpanSys",Gauge, float64(rtm.MSpanSys)},
		{"Mallocs", Gauge,float64(rtm.Mallocs)},
		{"NextGC", Gauge,float64(rtm.NextGC)},
		{"NumForcedGC",Gauge, float64(rtm.NumForcedGC)},
		{"NumGC",Gauge,float64(rtm.NumGC)},
		{"OtherSys",Gauge, float64(rtm.OtherSys)},
		{"PauseTotalNs", Gauge,float64(rtm.PauseTotalNs)},
		{"StackInuse",Gauge, float64(rtm.StackInuse)},
		{"StackSys", Gauge,float64(rtm.StackSys)},
		{"Sys", Gauge,float64(rtm.Sys)},
	}
	return m
}

func GetRandomValueMetric() Metric {
	return Metric{"RandomValue", Gauge, rand.Float64()}
}

func GetCounterMetric() Metric {
	PollCount++
	rand.Seed(time.Now().UnixNano())
	return Metric{"PollCount", Counter, float64(PollCount)}
}

func ResetCounterMetric() {
	PollCount = 0
}
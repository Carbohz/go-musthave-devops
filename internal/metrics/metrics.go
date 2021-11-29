package metrics

import (
	"math/rand"
	"runtime"
)

type Metric struct {
	Name     string
	Typename string
	Value    float64
}

func GetRuntimeMetrics() []Metric {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	m := []Metric{
		{"Alloc","gauge", float64(rtm.Alloc)},
		{"BuckHashSys","gauge", float64(rtm.BuckHashSys)},
		{"Frees", "gauge",float64(rtm.Frees)},
		{"GCCPUFraction","gauge", rtm.GCCPUFraction},
		{"GCSys","gauge", float64(rtm.GCSys)},
		{"HeapAlloc","gauge", float64(rtm.HeapAlloc)},
		{"HeapIdle","gauge", float64(rtm.HeapIdle)},
		{"HeapInuse", "gauge",float64(rtm.HeapInuse)},
		{"HeapObjects","gauge", float64(rtm.HeapObjects)},
		{"HeapReleased","gauge", float64(rtm.HeapReleased)},
		{"HeapSys","gauge", float64(rtm.HeapSys)},
		{"LastGC", "gauge",float64(rtm.LastGC)},
		{"Lookups","gauge", float64(rtm.Lookups)},
		{"MCacheInuse", "gauge",float64(rtm.MCacheInuse)},
		{"MCacheSys", "gauge",float64(rtm.MCacheSys)},
		{"MSpanInuse","gauge",float64(rtm.MSpanInuse)},
		{"MSpanSys","gauge", float64(rtm.MSpanSys)},
		{"Mallocs", "gauge",float64(rtm.Mallocs)},
		{"NextGC", "gauge",float64(rtm.NextGC)},
		{"NumForcedGC","gauge", float64(rtm.NumForcedGC)},
		{"NumGC","gauge",float64(rtm.NumGC)},
		{"OtherSys","gauge", float64(rtm.OtherSys)},
		{"PauseTotalNs", "gauge",float64(rtm.PauseTotalNs)},
		{"StackInuse","gauge", float64(rtm.StackInuse)},
		{"StackSys", "gauge",float64(rtm.StackSys)},
		{"Sys", "gauge",float64(rtm.Sys)},
	}
	return m
}

func GetRandomValueMetric() Metric {
	return Metric{"RandomValue", "gauge", rand.Float64()}
}

func GetCounterMetric(PollCount int64) Metric {
	return Metric{"PollCount", "counter", float64(PollCount)}
}
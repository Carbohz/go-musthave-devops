package agent

import (
	"github.com/Carbohz/go-musthave-devops/model"
	"math/rand"
	"runtime"
	"time"
)

func (agent *Agent) collectMetrics() {
	agent.metrics.memStats = collectMemStats()
	agent.metrics.randomValue = collectRandomValue()
	agent.metrics.pollCount.Value += 1
}

func collectMemStats() []model.GaugeMetric {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m := []model.GaugeMetric{
		{model.Common{"Alloc", model.Gauge}, float64(rtm.Alloc)},
		//{Common: model.Common{Name: "Alloc", Typename: model.Gauge}, Value: float64(rtm.Alloc)},
		{model.Common{Name: "BuckHashSys", Typename: model.Gauge}, float64(rtm.BuckHashSys)},
		{model.Common{Name: "Frees", Typename: model.Gauge},float64(rtm.Frees)},
		{model.Common{Name: "GCCPUFraction", Typename: model.Gauge}, rtm.GCCPUFraction},
		{model.Common{Name: "GCSys", Typename: model.Gauge}, float64(rtm.GCSys)},
		{model.Common{Name: "HeapAlloc", Typename: model.Gauge}, float64(rtm.HeapAlloc)},
		{model.Common{Name: "HeapIdle", Typename: model.Gauge}, float64(rtm.HeapIdle)},
		{model.Common{Name: "HeapInuse", Typename: model.Gauge},float64(rtm.HeapInuse)},
		{model.Common{Name: "HeapObjects", Typename: model.Gauge}, float64(rtm.HeapObjects)},
		{model.Common{Name: "HeapReleased", Typename: model.Gauge}, float64(rtm.HeapReleased)},
		{model.Common{Name: "HeapSys", Typename: model.Gauge}, float64(rtm.HeapSys)},
		{model.Common{Name: "LastGC", Typename: model.Gauge},float64(rtm.LastGC)},
		{model.Common{Name: "Lookups", Typename: model.Gauge}, float64(rtm.Lookups)},
		{model.Common{Name: "MCacheInuse", Typename: model.Gauge},float64(rtm.MCacheInuse)},
		{model.Common{Name: "MCacheSys", Typename: model.Gauge},float64(rtm.MCacheSys)},
		{model.Common{Name: "MSpanInuse", Typename: model.Gauge},float64(rtm.MSpanInuse)},
		{model.Common{Name: "MSpanSys", Typename: model.Gauge}, float64(rtm.MSpanSys)},
		{model.Common{Name: "Mallocs", Typename: model.Gauge},float64(rtm.Mallocs)},
		{model.Common{Name: "NextGC", Typename: model.Gauge},float64(rtm.NextGC)},
		{model.Common{Name: "NumForcedGC", Typename: model.Gauge}, float64(rtm.NumForcedGC)},
		{model.Common{Name: "NumGC", Typename: model.Gauge},float64(rtm.NumGC)},
		{model.Common{Name: "OtherSys", Typename: model.Gauge}, float64(rtm.OtherSys)},
		{model.Common{Name: "PauseTotalNs", Typename: model.Gauge},float64(rtm.PauseTotalNs)},
		{model.Common{Name: "StackInuse", Typename: model.Gauge}, float64(rtm.StackInuse)},
		{model.Common{Name: "StackSys", Typename: model.Gauge},float64(rtm.StackSys)},
		{model.Common{Name: "Sys", Typename: model.Gauge},float64(rtm.Sys)},
		{model.Common{Name: "TotalAlloc", Typename: model.Gauge}, float64(rtm.TotalAlloc)},
	}
	return m
}

func collectRandomValue() model.GaugeMetric {
	rand.Seed(time.Now().UnixNano())
	return model.GaugeMetric{Common: model.Common{Name: "RandomValue",  Typename: model.Gauge}, Value: rand.Float64()}
}


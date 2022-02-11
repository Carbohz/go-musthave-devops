package agent

import (
	"github.com/Carbohz/go-musthave-devops/model"
	"github.com/markphelps/optional"
	"math/rand"
	"runtime"
	"time"
)

func (agent *Agent) collectMetrics() {
	agent.metrics.memStats = collectMemStats()
	agent.metrics.randomValue = collectRandomValue()

	pollCount, _ := agent.metrics.pollCount.Delta.Get()
	agent.metrics.pollCount.Delta.Set(pollCount + 1)
}

func collectMemStats() []model.Metric {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m := []model.Metric{
		{Name: "Alloc", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.Alloc))},
		{Name: "BuckHashSys", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.BuckHashSys))},
		{Name: "Frees", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.Frees))},
		{Name: "GCCPUFraction", Type: model.KGauge, Value: optional.NewFloat64(rtm.GCCPUFraction)},
		{Name: "GCSys", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.GCSys))},
		{Name: "HeapAlloc", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.HeapAlloc))},
		{Name: "HeapIdle", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.HeapIdle))},
		{Name: "HeapInuse", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.HeapInuse))},
		{Name: "HeapObjects", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.HeapObjects))},
		{Name: "HeapReleased", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.HeapReleased))},
		{Name: "HeapSys", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.HeapSys))},
		{Name: "LastGC", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.LastGC))},
		{Name: "Lookups", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.Lookups))},
		{Name: "MCacheInuse", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.MCacheInuse))},
		{Name: "MCacheSys", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.MCacheSys))},
		{Name: "MSpanInuse", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.MSpanInuse))},
		{Name: "MSpanSys", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.MSpanSys))},
		{Name: "Mallocs", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.Mallocs))},
		{Name: "NextGC", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.NextGC))},
		{Name: "NumForcedGC", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.NumForcedGC))},
		{Name: "NumGC", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.NumGC))},
		{Name: "OtherSys", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.OtherSys))},
		{Name: "PauseTotalNs", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.PauseTotalNs))},
		{Name: "StackInuse", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.StackInuse))},
		{Name: "StackSys", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.StackSys))},
		{Name: "Sys", Type: model.KGauge,Value: optional.NewFloat64(float64(rtm.Sys))},
		{Name: "TotalAlloc", Type: model.KGauge, Value: optional.NewFloat64(float64(rtm.TotalAlloc))},
	}

	return m
}

func collectRandomValue() model.Metric {
	rand.Seed(time.Now().UnixNano())
	randomValue := optional.NewFloat64(rand.Float64())
	return model.Metric{Name: "RandomValue", Type: model.KGauge, Value: randomValue}
}

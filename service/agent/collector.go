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
	//agent.metrics.pollCount.Value += 1

	pollCount, _ := agent.metrics.pollCount.Delta.Get()
	agent.metrics.pollCount.Delta.Set(pollCount + 1)
}

func collectMemStats() []model.Metric {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	//m := []model.GaugeMetric{
	//	{Common: model.Common{Name: "Alloc", Typename: model.Gauge}, Value: float64(rtm.Alloc)},
	//	{Common: model.Common{Name: "BuckHashSys", Typename: model.Gauge}, Value: float64(rtm.BuckHashSys)},
	//	{Common: model.Common{Name: "Frees", Typename: model.Gauge},Value: float64(rtm.Frees)},
	//	{Common: model.Common{Name: "GCCPUFraction", Typename: model.Gauge}, Value: rtm.GCCPUFraction},
	//	{Common: model.Common{Name: "GCSys", Typename: model.Gauge}, Value: float64(rtm.GCSys)},
	//	{Common: model.Common{Name: "HeapAlloc", Typename: model.Gauge}, Value: float64(rtm.HeapAlloc)},
	//	{Common: model.Common{Name: "HeapIdle", Typename: model.Gauge}, Value: float64(rtm.HeapIdle)},
	//	{Common: model.Common{Name: "HeapInuse", Typename: model.Gauge},Value: float64(rtm.HeapInuse)},
	//	{Common: model.Common{Name: "HeapObjects", Typename: model.Gauge}, Value: float64(rtm.HeapObjects)},
	//	{Common: model.Common{Name: "HeapReleased", Typename: model.Gauge}, Value: float64(rtm.HeapReleased)},
	//	{Common: model.Common{Name: "HeapSys", Typename: model.Gauge}, Value: float64(rtm.HeapSys)},
	//	{Common: model.Common{Name: "LastGC", Typename: model.Gauge},Value: float64(rtm.LastGC)},
	//	{Common: model.Common{Name: "Lookups", Typename: model.Gauge}, Value: float64(rtm.Lookups)},
	//	{Common: model.Common{Name: "MCacheInuse", Typename: model.Gauge},Value: float64(rtm.MCacheInuse)},
	//	{Common: model.Common{Name: "MCacheSys", Typename: model.Gauge},Value: float64(rtm.MCacheSys)},
	//	{Common: model.Common{Name: "MSpanInuse", Typename: model.Gauge},Value: float64(rtm.MSpanInuse)},
	//	{Common: model.Common{Name: "MSpanSys", Typename: model.Gauge}, Value: float64(rtm.MSpanSys)},
	//	{Common: model.Common{Name: "Mallocs", Typename: model.Gauge},Value: float64(rtm.Mallocs)},
	//	{Common: model.Common{Name: "NextGC", Typename: model.Gauge},Value: float64(rtm.NextGC)},
	//	{Common: model.Common{Name: "NumForcedGC", Typename: model.Gauge}, Value: float64(rtm.NumForcedGC)},
	//	{Common: model.Common{Name: "NumGC", Typename: model.Gauge},Value: float64(rtm.NumGC)},
	//	{Common: model.Common{Name: "OtherSys", Typename: model.Gauge}, Value: float64(rtm.OtherSys)},
	//	{Common: model.Common{Name: "PauseTotalNs", Typename: model.Gauge},Value: float64(rtm.PauseTotalNs)},
	//	{Common: model.Common{Name: "StackInuse", Typename: model.Gauge}, Value: float64(rtm.StackInuse)},
	//	{Common: model.Common{Name: "StackSys", Typename: model.Gauge},Value: float64(rtm.StackSys)},
	//	{Common: model.Common{Name: "Sys", Typename: model.Gauge},Value: float64(rtm.Sys)},
	//	{Common: model.Common{Name: "TotalAlloc", Typename: model.Gauge}, Value: float64(rtm.TotalAlloc)},
	//}

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
	//return model.GaugeMetric{Common: model.Common{Name: "RandomValue",  Typename: model.Gauge}, Value: rand.Float64()}
	randomValue := optional.NewFloat64(rand.Float64())
	return model.Metric{Name: "RandomValue", Type: model.KGauge, Value: randomValue}
}


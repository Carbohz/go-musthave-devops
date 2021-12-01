package main

import (
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"github.com/Carbohz/go-musthave-devops/internal/sender"
	"net/http"
	"time"
)

const (
	pollInterval = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	RunAgent()
}

func RunAgent() {
	var PollCount int64
	var runtimeMetrics []metrics.Metric
	var randomValueMetric metrics.Metric
	var counterMetric metrics.Metric

	client := http.Client{Timeout: 2 * time.Second}
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval )
	for {
		select {
			case <-pollTicker.C:
				PollCount++
				runtimeMetrics = metrics.GetRuntimeMetrics()
				randomValueMetric = metrics.GetRandomValueMetric()
				counterMetric = metrics.GetCounterMetric(PollCount)
			case <-reportTicker.C:
				for _, m := range runtimeMetrics {
					sender.Send(&client, m)
				}
				sender.Send(&client, randomValueMetric)
				sender.Send(&client, counterMetric)
				PollCount = 0
		}
	}
}

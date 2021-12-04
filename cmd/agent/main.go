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
	var runtimeMetrics []metrics.GaugeMetric
	var randomValueMetric metrics.GaugeMetric
	var counterMetric metrics.CounterMetric

	client := http.Client{Timeout: 2 * time.Second}

	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval )
	for {
		select {
			case <-pollTicker.C:
				metrics.IncrementCounterMetric()
				runtimeMetrics = metrics.GetRuntimeMetrics()
				randomValueMetric = metrics.GetRandomValueMetric()
				counterMetric = metrics.GetCounterMetric()
			case <-reportTicker.C:
				for _, m := range runtimeMetrics {
					sender.SendGaugeMetric(&client, m)
				}
				sender.SendGaugeMetric(&client, randomValueMetric)
				sender.SendCounterMetric(&client, counterMetric)
				metrics.ResetCounterMetric()
		}
	}
}

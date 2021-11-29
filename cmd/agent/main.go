package main

import (
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"github.com/Carbohz/go-musthave-devops/internal/sender"
	"net/http"
	"time"
)

const (
	pollInterval = 2
	reportInterval = 10
)

func main() {
	RunAgent()
}

func RunAgent() {
	var PollCount int64

	client := http.Client{Timeout: 2 * time.Second}
	ticker := time.NewTicker(pollInterval * time.Second)
	for {
		<-ticker.C
		PollCount += 1
		runtimeMetrics := metrics.GetRuntimeMetrics()
		randomValueMetric := metrics.GetRandomValueMetric()
		counterMetric := metrics.GetCounterMetric(PollCount)

		if pollInterval * PollCount == reportInterval {
			for _, m := range runtimeMetrics {
				sender.Send(&client, m)
			}
			sender.Send(&client, randomValueMetric)
			sender.Send(&client, counterMetric)
			PollCount = 0
		}
	}
}

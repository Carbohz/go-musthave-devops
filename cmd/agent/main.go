package main

import (
	"github.com/Carbohz/go-musthave-devops/internal/common"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Carbohz/go-musthave-devops/internal/agent"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"github.com/Carbohz/go-musthave-devops/internal/sender"
)

func main() {
	cfg := agent.CreateConfig()
	exitChan := make(chan int, 1)
	go common.AwaitInterruptSignal(exitChan)
	go RunAgent(cfg)
	exitCode := <-exitChan
	os.Exit(exitCode)
}

func RunAgent(cfg agent.Config) {
	var runtimeMetrics []metrics.GaugeMetric
	var randomValueMetric metrics.GaugeMetric
	var pollCountMetric metrics.CounterMetric

	client := http.Client{Timeout: 2 * time.Second}

	pollTicker := time.NewTicker(cfg.PollInterval)
	reportTicker := time.NewTicker(cfg.ReportInterval)
	for {
		select {
		case <-pollTicker.C:
			log.Println("Collecting Metrics")
			metrics.IncrementPollCountMetric()
			runtimeMetrics = metrics.GetRuntimeMetrics()
			randomValueMetric = metrics.GetRandomValueMetric()
			pollCountMetric = metrics.GetPollCountMetric()
		case <-reportTicker.C:
			log.Println("Sending Metrics")
			for _, m := range runtimeMetrics {
				//sender.SendGaugeMetric(&client, m, cfg.Address)
				sender.SendGaugeMetricJSON(&client, m, cfg)
			}
			//sender.SendGaugeMetric(&client, randomValueMetric, cfg.Address)
			sender.SendGaugeMetricJSON(&client, randomValueMetric, cfg)

			//sender.SendCounterMetric(&client, pollCountMetric, cfg.Address)
			sender.SendCounterMetricJSON(&client, pollCountMetric, cfg)

			//sender.SendMetricsJSON(&client, runtimeMetrics, randomValueMetric, pollCountMetric, cfg)

			metrics.ResetPollCountMetric()
		}
	}
}

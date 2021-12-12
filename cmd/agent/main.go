package main

import (
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"github.com/Carbohz/go-musthave-devops/internal/sender"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string         `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

const (
	defaultAddress = "127.0.0.1:8080"
	defaultPollInterval = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Address == "" {
		cfg.Address = defaultAddress
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = defaultReportInterval
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = defaultPollInterval
	}

	RunAgent(cfg)
}

func RunAgent(cfg Config) {
	var runtimeMetrics []metrics.GaugeMetric
	var randomValueMetric metrics.GaugeMetric
	var pollCountMetric metrics.CounterMetric

	client := http.Client{Timeout: 2 * time.Second}

	pollTicker := time.NewTicker(cfg.PollInterval)
	reportTicker := time.NewTicker(cfg.ReportInterval)
	for {
		select {
			case <-pollTicker.C:
				metrics.IncrementPollCountMetric()
				runtimeMetrics = metrics.GetRuntimeMetrics()
				randomValueMetric = metrics.GetRandomValueMetric()
				pollCountMetric = metrics.GetPollCountMetric()
			case <-reportTicker.C:
				for _, m := range runtimeMetrics {
					sender.SendGaugeMetric(&client, m, cfg.Address)
				}
				sender.SendGaugeMetric(&client, randomValueMetric, cfg.Address)
				sender.SendCounterMetric(&client, pollCountMetric, cfg.Address)
				metrics.ResetPollCountMetric()
		}
	}
}

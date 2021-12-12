package main

import (
	"flag"
	"github.com/Carbohz/go-musthave-devops/internal/metrics"
	"github.com/Carbohz/go-musthave-devops/internal/sender"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string         `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
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

	addressFlagPtr := flag.String("a", defaultAddress, "set server's address where you want to send metrics")
	pollIntervalFlagPtr := flag.Duration("p", defaultPollInterval, "set metrics poll interval")
	reportIntervalFlagPtr := flag.Duration("r", defaultReportInterval, "set metrics report interval")

	flag.Parse()

	_, isSet := os.LookupEnv("ADDRESS")
	if !isSet {
		if addressFlagPtr != nil {
			cfg.Address = *addressFlagPtr
		} else {
			cfg.Address = defaultAddress
		}
	}

	_, isSet = os.LookupEnv("POLL_INTERVAL")
	if !isSet {
		if pollIntervalFlagPtr != nil {
			cfg.PollInterval = *pollIntervalFlagPtr
		} else {
			cfg.PollInterval = defaultPollInterval
		}
	}

	_, isSet = os.LookupEnv("REPORT_INTERVAL")
	if !isSet {
		if reportIntervalFlagPtr != nil {
			cfg.ReportInterval = *reportIntervalFlagPtr
		} else {
			cfg.ReportInterval = defaultReportInterval
		}
	}

	//if cfg.Address == "" {
	//	cfg.Address = defaultAddress
	//}
	//
	//if cfg.ReportInterval == 0 {
	//	cfg.ReportInterval = defaultReportInterval
	//}
	//
	//if cfg.PollInterval == 0 {
	//	cfg.PollInterval = defaultPollInterval
	//}

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
				log.Println("Collecting Metrics")
				metrics.IncrementPollCountMetric()
				runtimeMetrics = metrics.GetRuntimeMetrics()
				randomValueMetric = metrics.GetRandomValueMetric()
				pollCountMetric = metrics.GetPollCountMetric()
			case <-reportTicker.C:
				log.Println("Sending Metrics")
				for _, m := range runtimeMetrics {
					sender.SendGaugeMetric(&client, m, cfg.Address)
				}
				sender.SendGaugeMetric(&client, randomValueMetric, cfg.Address)
				sender.SendCounterMetric(&client, pollCountMetric, cfg.Address)
				metrics.ResetPollCountMetric()
		}
	}
}

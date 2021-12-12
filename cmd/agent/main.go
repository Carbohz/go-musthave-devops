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
	cfg := CreateConfig()
	RunAgent(cfg)
}

func CreateConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Agent is running with environment variables: %+v", cfg)

	addressFlagPtr := flag.String("a", defaultAddress, "set server's address where you want to send metrics")
	pollIntervalFlagPtr := flag.Duration("p", defaultPollInterval, "set metrics poll interval")
	reportIntervalFlagPtr := flag.Duration("r", defaultReportInterval, "set metrics report interval")
	flag.Parse()
	log.Printf("Agent is running with command line flags: Address %v, Poll Interval %v, Report Interval %v",
		*addressFlagPtr, *pollIntervalFlagPtr, *reportIntervalFlagPtr)

	_, isSet := os.LookupEnv("ADDRESS")
	if !isSet {
		cfg.Address = *addressFlagPtr
	}

	_, isSet = os.LookupEnv("POLL_INTERVAL")
	if !isSet {
		cfg.PollInterval = *pollIntervalFlagPtr
	}

	_, isSet = os.LookupEnv("REPORT_INTERVAL")
	if !isSet {
		cfg.ReportInterval = *reportIntervalFlagPtr
	}
	
	return cfg
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

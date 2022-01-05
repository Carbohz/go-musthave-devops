package agent

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"time"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	Key            string        `env:"KEY"`
}

const (
	defaultAddress        = "127.0.0.1:8080"
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
	defaultKeyHash        = "abracadabra"
)

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
	keyHashFlagPtr := flag.String("k", defaultKeyHash, "enter key to compute hash for safe data sending")
	flag.Parse()
	log.Printf("Agent is running with command line flags: Address %v, Poll Interval %v, Report Interval %v, Key %v",
		*addressFlagPtr, *pollIntervalFlagPtr, *reportIntervalFlagPtr, *keyHashFlagPtr)

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

	_, isSet = os.LookupEnv("KEY")
	if !isSet {
		cfg.Key = *keyHashFlagPtr
	}

	return cfg
}

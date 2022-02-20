package configagent

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type AgentConfig struct {
	Address        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	Key            string        `env:"KEY"`
}

const (
	defaultAddress        = "127.0.0.1:8080"
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
	defaultKeyHash        = ""
)

func NewAgentConfig() (AgentConfig, error) {
	var cfg AgentConfig

	err := env.Parse(&cfg)
	if err != nil {
		return AgentConfig{}, fmt.Errorf("agent config Ctor error : %w", err)
	}

	log.Printf("Agent is running with environment variables: %+v", cfg)

	addressFlagPtr := flag.String("a", defaultAddress, "set server's address where you want to send metrics")
	pollIntervalFlagPtr := flag.Duration("p", defaultPollInterval, "set metrics poll interval")
	reportIntervalFlagPtr := flag.Duration("r", defaultReportInterval, "set metrics report interval")
	keyHashFlagPtr := flag.String("k", defaultKeyHash, "enter key to compute hash for safe data sending")
	flag.Parse()
	log.Printf("Agent is running with command line flags: Address %v, Poll Interval %v, Report Interval %v, Key %v",
		*addressFlagPtr, *pollIntervalFlagPtr, *reportIntervalFlagPtr, *keyHashFlagPtr)

	if _, isSet := os.LookupEnv("ADDRESS"); !isSet {
		cfg.Address = *addressFlagPtr
	}

	if _, isSet := os.LookupEnv("POLL_INTERVAL"); !isSet {
		cfg.PollInterval = *pollIntervalFlagPtr
	}

	if _, isSet := os.LookupEnv("REPORT_INTERVAL"); !isSet {
		cfg.ReportInterval = *reportIntervalFlagPtr
	}

	if _, isSet := os.LookupEnv("KEY"); !isSet {
		cfg.Key = *keyHashFlagPtr
	}

	log.Printf("Final Agent configuration: %+v", cfg)

	return cfg, nil
}

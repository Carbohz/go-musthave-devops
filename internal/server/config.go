package server

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"time"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
}

const (
	defaultAddress       = "127.0.0.1:8080"
	defaultStoreInterval = 300 * time.Second
	defaultStoreFile     = "/tmp/devops-metrics-db.json"
	defaultRestore       = true
	defaultKeyHash       = ""
)

func CreateConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server is running with environment variables: %+v", cfg)

	addressFlagPtr := flag.String("a", defaultAddress, "set address of server")
	storeIntervalFlagPtr := flag.Duration("i", defaultStoreInterval, "set server's metrics store interval")
	storeFileFlagPtr := flag.String("f", defaultStoreFile, "set file where metrics are stored")
	restoreFlagPtr := flag.Bool("r", defaultRestore, "choose whether to restore server metrics from file")
	keyHashFlagPtr := flag.String("k", defaultKeyHash, "enter key to compute hash for safe data sending")
	flag.Parse()
	log.Printf("Server is running with command line flags: Address %v, Store Interval %v, Store File %v, Restore %v, Key %v",
		*addressFlagPtr, *storeIntervalFlagPtr, *storeFileFlagPtr, *restoreFlagPtr, *keyHashFlagPtr)

	_, isSet := os.LookupEnv("ADDRESS")
	if !isSet {
		cfg.Address = *addressFlagPtr
	}

	_, isSet = os.LookupEnv("STORE_INTERVAL")
	if !isSet {
		cfg.StoreInterval = *storeIntervalFlagPtr
	}

	_, isSet = os.LookupEnv("STORE_FILE")
	if !isSet {
		cfg.StoreFile = *storeFileFlagPtr
	}

	_, isSet = os.LookupEnv("RESTORE")
	if !isSet {
		cfg.Restore = *restoreFlagPtr
	}

	_, isSet = os.LookupEnv("KEY")
	if !isSet {
		cfg.Key = *keyHashFlagPtr
	}

	return cfg
}

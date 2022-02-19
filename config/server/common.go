package configsrv

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

type CommonConfig struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DBPath        string        `env:"DATABASE_DSN"`
}

const (
	defaultAddress       = "127.0.0.1:8080"
	defaultStoreInterval = 300 * time.Second
	defaultStoreFile     = "D:\\Go\\yandex-praktikum\\increments\\go-musthave-devops\\tmp\\devops-metrics-db.json"
	defaultRestore       = true
	defaultKeyHash       = ""
	defaultDBPath        = ""
)

func NewCommonConfig() (CommonConfig, error) {
	var cfg CommonConfig

	if err := env.Parse(&cfg); err != nil {
		return CommonConfig{}, fmt.Errorf("common config Ctor error : %w", err)
	}
	log.Printf("Server is running with environment variables: %+v", cfg)

	addressFlagPtr := flag.String("a", defaultAddress, "set address of server")
	storeIntervalFlagPtr := flag.Duration("i", defaultStoreInterval, "set server's metrics store interval")
	storeFileFlagPtr := flag.String("f", defaultStoreFile, "set file where metrics are stored")
	restoreFlagPtr := flag.Bool("r", defaultRestore, "choose whether to restore server metrics from file")
	keyHashFlagPtr := flag.String("k", defaultKeyHash, "enter key to compute hash for safe data sending")
	dbPathPtr := flag.String("d", defaultDBPath, "set address of db to connect")
	flag.Parse()
	log.Printf("Server is running with command line flags: Address %s, Store Interval %v, Store File %s, Restore %v, Key %s, DB: %s",
		*addressFlagPtr, *storeIntervalFlagPtr, *storeFileFlagPtr, *restoreFlagPtr, *keyHashFlagPtr, *dbPathPtr)

	if _, isSet := os.LookupEnv("ADDRESS"); !isSet {
		cfg.Address = *addressFlagPtr
	}

	if _, isSet := os.LookupEnv("STORE_INTERVAL"); !isSet {
		cfg.StoreInterval = *storeIntervalFlagPtr
	}

	if _, isSet := os.LookupEnv("STORE_FILE"); !isSet {
		cfg.StoreFile = *storeFileFlagPtr
	}

	if _, isSet := os.LookupEnv("RESTORE"); !isSet {
		cfg.Restore = *restoreFlagPtr
	}

	if _, isSet := os.LookupEnv("KEY"); !isSet {
		cfg.Key = *keyHashFlagPtr
	}

	if _, isSet := os.LookupEnv("DATABASE_DSN"); !isSet {
		cfg.DBPath = *dbPathPtr
	}

	log.Printf("Final server configuration: %+v", cfg)

	return cfg, nil
}

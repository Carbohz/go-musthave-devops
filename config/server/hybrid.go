package configsrv

import (
	"time"
)

type HybridStorageConfig struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	DBPath        string        `env:"DATABASE_DSN"`
}

func NewHybridStorageConfig(config CommonConfig) HybridStorageConfig {
	hybridStorageConfig := HybridStorageConfig{
		StoreInterval: config.StoreInterval,
		StoreFile: config.StoreFile,
		Restore: config.Restore,
		DBPath: config.DBPath,
	}

	return hybridStorageConfig
}
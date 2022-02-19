package configsrv

import (
	"time"
)

type HybridStorageConfig struct {
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	DBPath        string
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
package configsrv

import "time"

type ServerConfig struct {
	Address       string
	StoreInterval time.Duration
	Key           string
}

func NewServerConfig(config CommonConfig) ServerConfig {
	serverConfig := ServerConfig{
		Address: config.Address,
		StoreInterval: config.StoreInterval,
		Key: config.Key,
	}

	return serverConfig
}

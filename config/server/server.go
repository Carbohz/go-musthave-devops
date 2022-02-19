package configsrv

import "time"

type ServerConfig struct {
	Address       string
	StoreInterval time.Duration
	//StoreFile     string
	//Restore       bool
	Key           string
	//DBPath        string
}

func NewServerConfig(config CommonConfig) ServerConfig {
	serverConfig := ServerConfig{
		Address: config.Address,
		StoreInterval: config.StoreInterval,
		//StoreFile: config.StoreFile,
		//Restore: config.Restore,
		Key: config.Key,
		//DBPath: config.DBPath,
	}

	return serverConfig
}

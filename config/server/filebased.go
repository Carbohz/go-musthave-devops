package configsrv

type FileBasedStorageConfig struct {
	//StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func NewFileBasedStorageConfig(config CommonConfig) FileBasedStorageConfig {
	fileBasedConfig := FileBasedStorageConfig{
		//StoreInterval: config.StoreInterval,
		StoreFile: config.StoreFile,
		Restore: config.Restore,
	}

	return fileBasedConfig
}
package configsrv

type FileBasedStorageConfig struct {
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func NewFileBasedStorageConfig(config CommonConfig) FileBasedStorageConfig {
	fileBasedConfig := FileBasedStorageConfig{
		StoreFile: config.StoreFile,
		Restore: config.Restore,
	}

	return fileBasedConfig
}
package configsrv

type FileBasedStorageConfig struct {
	StoreFile     string
	Restore       bool
}

func NewFileBasedStorageConfig(config CommonConfig) FileBasedStorageConfig {
	fileBasedConfig := FileBasedStorageConfig{
		StoreFile: config.StoreFile,
		Restore: config.Restore,
	}

	return fileBasedConfig
}
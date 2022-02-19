package configsrv

type DatabaseStorageConfig struct {
	DBPath string
}

func NewDatabaseConfig(config CommonConfig) DatabaseStorageConfig {
	databaseConfig := DatabaseStorageConfig{
		DBPath: config.DBPath,
	}

	return databaseConfig
}

package config

// DBConfig config of database
type DBConfig struct {
	DBPath string `toml:"DBPath"`
}

func defaultDBConfig() DBConfig {
	return DBConfig{
		DBPath: "./dbstorage",
	}
}

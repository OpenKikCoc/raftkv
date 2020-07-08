package config

// GRPCConfig the config of grpc service
type GRPCConfig struct {
	HTTPPort string `toml:"HttpPort"`
	CertPath string `toml:"CertPath"`
	KeyPath  string `toml:"KeyPath"`
	Host     string `toml:"Host"`
}

func defaultGRPCConfig() GRPCConfig {
	return GRPCConfig{
		HTTPPort: "9500",
		CertPath: "server.crt",
		KeyPath:  "server.key",
		Host:     "binacs.cn",
	}
}

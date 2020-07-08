package config

// WebConfig the config of web service
type WebConfig struct {
	HTTPPort   string `toml:"HttpPort"`
	HTTPSPort  string `toml:"HttpsPort"`
	StaticPath string `toml:"StaticPath"` // end with '/'
	CertPath   string `toml:"CertPath"`
	KeyPath    string `toml:"KeyPath"`
	Host       string `toml:"Host"`
}

func defaultWebConfig() WebConfig {
	return WebConfig{
		HTTPPort:   "80",
		HTTPSPort:  "443",
		StaticPath: "./",
		CertPath:   "server.crt",
		KeyPath:    "server.key",
		Host:       "binacs.cn",
	}
}

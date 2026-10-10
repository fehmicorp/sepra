package main

type Config struct {
	Port            int    `yaml:"port" json:"port" env:"APP_PORT" default:"9090"`
	AccessKeyID     string `yaml:"access_key_id" json:"access_key_id" env:"APP_ACCESS_KEY_ID" default:"minioadmin"`
	SecretAccessKey string `yaml:"secret_access_key" json:"secret_access_key" env:"APP_SECRET_ACCESS_KEY" default:"minioadmin"`
	UseSSL          bool   `yaml:"use_ssl" json:"use_ssl" env:"APP_USE_SSL" default:"false"`
	MaxConnections  int    `yaml:"max_connections" json:"max_connections" env:"APP_MAX_CONNECTIONS" default:"10"`
	DebugMode       bool   `yaml:"debug_mode" json:"debug_mode" env:"APP_DEBUG_MODE" default:"false"`
	Timeout         string `yaml:"timeout" json:"timeout" env:"APP_TIMEOUT" default:"30s"`
}

package main

type Config struct {
	Port              int    `yaml:"port" json:"port" env:"APP_PORT" envDefault:"9098"`
	AccessKeyID       string `yaml:"access_key_id" json:"access_key_id" env:"APP_ACCESS_KEY_ID" envDefault:"minioadmin"`
	SecretAccessKey   string `yaml:"secret_access_key" json:"secret_access_key" env:"APP_SECRET_ACCESS_KEY" envDefault:"minioadmin"`
	SecretAccessToken string `yaml:"secret_access_token" json:"secret_access_token" env:"APP_SECRET_ACCESS_TOKEN" envDefault:"minioadmin"`
	UseSSL            bool   `yaml:"use_ssl" json:"use_ssl" env:"APP_USE_SSL" envDefault:"false"`
	MaxConnections    int    `yaml:"max_connections" json:"max_connections" env:"APP_MAX_CONNECTIONS" envDefault:"10"`
	DebugMode         bool   `yaml:"debug_mode" json:"debug_mode" env:"APP_DEBUG_MODE" envDefault:"false"`
	Timeout           string `yaml:"timeout" json:"timeout" env:"APP_TIMEOUT" envDefault:"30s"`
	Region            string `yaml:"region" json:"region" env:"APP_REGION" envDefault:"us-east-1"`
	LogLevel          string `yaml:"log_level" json:"log_level" env:"APP_LOG_LEVEL" envDefault:"info"`
	BackupDir         string `yaml:"backup_dir" json:"backup_dir" env:"APP_BACKUP_DIR" envDefault:"/backup"`
	BackupEnabled     bool   `yaml:"backup_enabled" json:"backup_enabled" env:"APP_BACKUP_ENABLED" envDefault:"true"`
	BackupInterval    int    `yaml:"backup_interval" json:"backup_interval" env:"APP_BACKUP_INTERVAL" envDefault:"86400"`
	DataDir           string `yaml:"data_dir" json:"data_dir" env:"APP_DATA_DIR" envDefault:"/data"`
}

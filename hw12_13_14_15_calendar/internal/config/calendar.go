package config

import "time"

type Config struct {
	Logger   LoggerConf   `mapstructure:"logger"`
	Server   ServerConf   `mapstructure:"server"`
	Storage  StorageConf  `mapstructure:"storage"`
	Database DatabaseConf `mapstructure:"database"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
}

type ServerConf struct {
	Address string `mapstructure:"address"`
}

type StorageConf struct {
	Type string `mapstructure:"type"`
}

type DatabaseConf struct {
	DSN             string        `mapstructure:"dsn"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

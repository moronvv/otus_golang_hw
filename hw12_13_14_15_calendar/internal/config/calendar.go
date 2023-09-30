package config

import "time"

type Config struct {
	Logger   *LoggerConf
	Server   *ServerConf
	Storage  *StorageConf
	Database *DatabaseConf
}

type LoggerConf struct {
	Level string
}

type ServerConf struct {
	Address string
}

type StorageConf struct {
	Type string
}

type DatabaseConf struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewConfig(
	loggerConf *LoggerConf,
	serverConf *ServerConf,
	storageConf *StorageConf,
	databaseConf *DatabaseConf,
) *Config {
	return &Config{
		Logger:   loggerConf,
		Server:   serverConf,
		Storage:  storageConf,
		Database: databaseConf,
	}
}

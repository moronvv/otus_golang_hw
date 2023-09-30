package config

type Config struct {
	Logger  LoggerConf
	Server  ServerConf
	Storage StorageConf
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

func NewConfig(
	loggerConf LoggerConf,
	serverConf ServerConf,
	storageConf StorageConf,
) *Config {
	return &Config{
		Logger:  loggerConf,
		Server:  serverConf,
		Storage: storageConf,
	}
}

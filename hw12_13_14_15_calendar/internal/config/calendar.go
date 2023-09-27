package config

type Config struct {
	Logger LoggerConf
	Server ServerConf
}

type LoggerConf struct {
	Level string
}

type ServerConf struct {
	Address string
}

func NewConfig(loggerConf LoggerConf, serverConf ServerConf) *Config {
	return &Config{
		Logger: loggerConf,
		Server: serverConf,
	}
}

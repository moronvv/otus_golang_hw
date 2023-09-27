package config

type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func NewConfig(logLvl string) *Config {
	return &Config{
		Logger: LoggerConf{
			Level: logLvl,
		},
	}
}

package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func NewConfig(cfgFile string) (*Config, error) {
	if err := initConfig(cfgFile); err != nil {
		return nil, err
	}

	return &Config{
		LoggerConf{
			Level: viper.GetString("logger.level"),
		},
	}, nil
}

func initConfig(cfgFile string) error {
	if cfgFile == "" {
		return fmt.Errorf("config file path not set")
	}
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	return viper.ReadInConfig()
}

// TODO

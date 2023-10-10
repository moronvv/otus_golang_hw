package main

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
)

func initConfig(cfgFile string) (*config.Config, error) {
	viper.SetConfigFile(cfgFile)

	viper.SetEnvPrefix("CALENDAR")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

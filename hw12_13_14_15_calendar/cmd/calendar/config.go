package main

import (
	"github.com/spf13/viper"

	"github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
)

func initConfig(cfgFile string) (*config.Config, error) {
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return config.NewConfig(
		&config.LoggerConf{
			Level: viper.GetString("logger.level"),
		},
		&config.ServerConf{
			Address: viper.GetString("server.address"),
		},
		&config.StorageConf{
			Type: viper.GetString("storage.type"),
		},
		&config.DatabaseConf{
			DSN:             viper.GetString("dsn"),
			MaxOpenConns:    viper.GetInt("max_open_conns"),
			MaxIdleConns:    viper.GetInt("max_idle_conns"),
			ConnMaxLifetime: viper.GetDuration("conn_max_lifetime"),
		},
	), nil
}

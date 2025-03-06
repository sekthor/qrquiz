package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database string
	Listen   string
	Loglevel string
}

func ReadConfig() (*Config, error) {
	var config Config
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.BindEnv("Database")
	viper.BindEnv("Listen")
	viper.BindEnv("Loglevel")

	if err := viper.Unmarshal(&config); err != nil {
		return &config, err
	}

	return &config, nil
}

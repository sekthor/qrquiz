package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database string
	Listen   string
	Loglevel string
	Otlp     struct {
		Enabled  bool
		Endpoint string
		Protocol string
		Interval int
		Insecure bool
	}
}

func ReadConfig() (*Config, error) {
	var config Config
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("Listen", ":8080")
	viper.SetDefault("Otlp.Endpoint", "localhost:4317")
	viper.SetDefault("Otlp.Interval", 10)
	viper.SetDefault("Otlp.Protocol", "grpc")
	viper.SetDefault("Otlp.Insecure", true)
	viper.BindEnv("Otlp.Enabled")
	viper.BindEnv("Otlp.Insecure")
	viper.BindEnv("Database")
	viper.BindEnv("Listen")
	viper.BindEnv("Loglevel")

	if err := viper.Unmarshal(&config); err != nil {
		return &config, err
	}

	return &config, nil
}

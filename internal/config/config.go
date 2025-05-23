package config

import (
	"strings"

	"github.com/nrednav/cuid2"
	"github.com/sirupsen/logrus"
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
	StaticCacheMaxAge int
	Contact           struct {
		Enabled      bool
		GeneralName  string
		GeneralEmail string
		AbuseName    string
		AbuseEmail   string
	}
	Admin struct {
		Password string
		User     string
		Disabled bool
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
	viper.BindEnv("StaticCacheMaxAge")
	viper.BindEnv("Contact.Enabled")
	viper.BindEnv("Contact.AbuseEmail")
	viper.BindEnv("Contact.AbuseName")
	viper.BindEnv("Contact.GeneralEmail")
	viper.BindEnv("Contact.GeneralName")
	viper.BindEnv("Admin.User")
	viper.BindEnv("Admin.Password")
	viper.BindEnv("Admin.Disabled")
	viper.SetDefault("Admin.User", "admin")
	viper.SetDefault("Admin.Password", cuid2.Generate()) // prevent misconfig empty password -> random password

	if err := viper.Unmarshal(&config); err != nil {
		return &config, err
	}

	return &config, nil
}

func (c Config) GetLoglevel() logrus.Level {
	switch c.Loglevel {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

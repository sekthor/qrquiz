package main

import (
	"context"
	"log"

	"github.com/sekthor/qrquiz/internal/config"
	"github.com/sekthor/qrquiz/internal/server"
	"github.com/sekthor/qrquiz/internal/telemetry"
	"github.com/sirupsen/logrus"
)

func main() {
	config, err := config.ReadConfig()
	if err != nil {
		logrus.WithField("error", err).Fatal("could not read config")
	}

	logrus.SetLevel(config.GetLoglevel())

	if config.Otlp.Enabled {
		shutdown, err := telemetry.SetUpTelemetry(context.Background(), config, "qrQuiz")
		if err != nil {
			log.Fatal(err)
		}
		defer shutdown(context.Background())
		logrus.Info("set up opentelemetry")
	}

	server := server.Server{}
	if err := server.Run(config); err != nil {
		log.Fatal(err)
	}
}

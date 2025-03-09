package main

import (
	"context"
	"log"
	"os"

	"github.com/sekthor/qrquiz/internal/config"
	"github.com/sekthor/qrquiz/internal/server"
	"github.com/sekthor/qrquiz/internal/telemetry"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

func main() {
	config, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if config.Otlp.Enabled {
		shutdown, err := telemetry.SetUpTelemetry(context.Background(), config, "qrQuiz")
		if err != nil {
			log.Fatal(err)
		}
		defer shutdown(context.Background())
	}

	// TODO: configre
	logrus.SetLevel(logrus.InfoLevel)

	server := server.Server{}
	if err := server.Run(config); err != nil {
		log.Fatal(err)
	}
}

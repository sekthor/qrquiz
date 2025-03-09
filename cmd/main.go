package main

import (
	"context"
	"log"

	"github.com/sekthor/qrquiz/internal/config"
	"github.com/sekthor/qrquiz/internal/server"
	"github.com/sekthor/qrquiz/internal/telemetry"
)

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

	server := server.Server{}
	if err := server.Run(config); err != nil {
		log.Fatal(err)
	}
}

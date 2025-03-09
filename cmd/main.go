package main

import (
	"log"
	"os"

	"github.com/sekthor/qrquiz/internal/config"
	"github.com/sekthor/qrquiz/internal/server"
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

	// TODO: configre
	logrus.SetLevel(logrus.InfoLevel)

	server := server.Server{}
	if err := server.Run(config); err != nil {
		log.Fatal(err)
	}
}

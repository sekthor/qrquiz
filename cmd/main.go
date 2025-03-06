package main

import (
	"log"

	"github.com/sekthor/qrquiz/internal/config"
	"github.com/sekthor/qrquiz/internal/server"
)

func main() {
	config, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	server := server.Server{}
	if err := server.Run(config); err != nil {
		log.Fatal(err)
	}
}

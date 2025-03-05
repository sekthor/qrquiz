package main

import (
	"log"

	"github.com/sekthor/qrquiz/internal/server"
)

func main() {
	server := server.Server{}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

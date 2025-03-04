package main

import (
	"log"

	"github.com/sekthor/puzzleinvite/internal/server"
)

func main() {
	server := server.Server{}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

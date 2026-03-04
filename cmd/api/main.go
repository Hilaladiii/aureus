package main

import (
	"log"

	"otoflash/pkg/config"
)

func main() {
	// will be implemented soon
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config", err)
	}
}

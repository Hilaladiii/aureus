package main

import (
	"log"

	di "github.com/Hilaladiii/aureus/di"
)

func main() {
	app, err := di.InitializeApp()
	if err != nil {
		log.Fatalf("failed initialized server: %v", err)
	}

	err = app.Listen(":8000")
	if err != nil {
		log.Fatalf("Server crash: %v", err)
	}
}

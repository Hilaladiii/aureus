package main

import (
	"context"
	"log"

	"github.com/Hilaladiii/aureus/di"
)

func main() {
	app, err := di.InitializeApp()
	if err != nil {
		log.Fatalf("failed initialized server: %v", err)
	}

	ctx := context.Background()
	app.Worker.Start(ctx)

	err = app.Web.Listen(":8000")
	if err != nil {
		log.Fatalf("Server crash: %v", err)
	}
}

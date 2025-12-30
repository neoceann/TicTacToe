package main

import (
	"log"
	"os"

	"go.uber.org/fx"
	"tictactoe/internal/di"
)

func main() {
	app := fx.New(
		di.Module,
		
		fx.NopLogger,
	)
	
	if err := app.Err(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
		os.Exit(1)
	}
	
	log.Println("Application initialized successfully")

	app.Run()
}
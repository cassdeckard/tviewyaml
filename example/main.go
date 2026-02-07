package main

import (
	"log"

	"github.com/cassdeckard/tviewyaml"
)

func main() {
	// Create app from YAML config directory
	// This assumes you have a "config" directory with root.yaml and page YAML files
	app, err := tviewyaml.CreateApp("./config")
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Run the application
	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

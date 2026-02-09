package main

import (
	"log"

	"github.com/cassdeckard/tviewyaml"
)

func main() {
	app, pageErrors, err := tviewyaml.NewAppBuilder("./config").
		With(RegisterClock).
		Build()
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Ensure cleanup of background goroutines
	defer app.Stop()

	if len(pageErrors) > 0 {
		log.Printf("Warning: %d page(s) failed to load/build:", len(pageErrors))
		for _, pageErr := range pageErrors {
			log.Printf("  - %v", pageErr)
		}
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

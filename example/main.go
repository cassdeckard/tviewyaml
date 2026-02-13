package main

import (
	"log"

	"example/app"
)

func main() {
	appObj, pageErrors, err := app.Build("./config")
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Ensure cleanup of background goroutines
	defer appObj.Stop()

	if len(pageErrors) > 0 {
		log.Printf("Warning: %d page(s) failed to load/build:", len(pageErrors))
		for _, pageErr := range pageErrors {
			log.Printf("  - %v", pageErr)
		}
	}

	if err := appObj.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

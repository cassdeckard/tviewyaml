package main

import (
	"log"

	"github.com/cassdeckard/tviewyaml"
)

func main() {
	app, err := tviewyaml.NewAppBuilder("./config").
		RegisterTemplateFunctions(RegisterClockFunctions).
		Build()
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

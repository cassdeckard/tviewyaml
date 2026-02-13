package app

import (
	"github.com/cassdeckard/tviewyaml"
	"github.com/gdamore/tcell/v2"
)

// Build creates the example application with the given config path.
// Use this from the main binary (default screen).
func Build(configPath string) (*tviewyaml.Application, []error, error) {
	return build(configPath, nil)
}

// BuildWithScreen creates the example application with the given config path and screen.
// Use this from acceptance tests with a simulation screen.
func BuildWithScreen(configPath string, screen tcell.Screen) (*tviewyaml.Application, []error, error) {
	return build(configPath, screen)
}

func build(configPath string, screen tcell.Screen) (*tviewyaml.Application, []error, error) {
	b := tviewyaml.NewAppBuilder(configPath).
		With(RegisterClock).
		With(RegisterStateBinding).
		With(RegisterInputFieldLive).
		With(RegisterDynamicPages)
	if screen != nil {
		b = b.WithScreen(screen)
	}
	return b.Build()
}

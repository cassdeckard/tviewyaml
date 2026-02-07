package tviewyaml

import (
	"log"

	"github.com/cassdeckard/tviewyaml/builder"
	"github.com/cassdeckard/tviewyaml/config"
	"github.com/cassdeckard/tviewyaml/template"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateApp creates and configures a tview application from YAML configuration files
func CreateApp(configDir string) (*tview.Application, error) {
	// Initialize tview application
	app := tview.NewApplication()
	pages := tview.NewPages()

	// Create template context
	ctx := template.NewContext(app, pages)

	// Load configuration
	loader := config.NewLoader(configDir)
	appConfig, err := loader.LoadApp("app.yaml")
	if err != nil {
		return nil, err
	}

	// Validate app config
	validator := config.NewValidator()
	if err := validator.ValidateApp(appConfig); err != nil {
		return nil, err
	}

	// Create builder
	uiBuilder := builder.NewBuilder(ctx)

	// Build all pages from config
	for _, pageRef := range appConfig.Application.Root.Pages {
		pageConfig, err := loader.LoadPage(pageRef.Ref)
		if err != nil {
			log.Printf("Error loading page %s: %v", pageRef.Name, err)
			continue
		}

		// Validate page config
		if err := validator.ValidatePage(pageConfig); err != nil {
			log.Printf("Invalid page config %s: %v", pageRef.Name, err)
			continue
		}

		pagePrimitive, err := uiBuilder.BuildFromConfig(pageConfig)
		if err != nil {
			log.Printf("Error building page %s: %v", pageRef.Name, err)
			continue
		}

		// Add to pages
		visible := pageRef.Name == "main"
		pages.AddPage(pageRef.Name, pagePrimitive, true, visible)
	}

	// Apply global keyboard shortcuts from YAML
	if len(appConfig.Application.GlobalKeyBindings) > 0 {
		executor := template.NewExecutor(ctx)
		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			for _, binding := range appConfig.Application.GlobalKeyBindings {
				if template.MatchesKeyBinding(event, binding) {
					callback, err := executor.ExecuteCallback(binding.Action)
					if err == nil {
						callback()
						return nil
					}
				}
			}
			return event
		})
	}

	// Apply mouse setting (default to true if not specified)
	enableMouse := true
	if appConfig.Application.EnableMouse {
		enableMouse = appConfig.Application.EnableMouse
	}

	return app.SetRoot(pages, true).EnableMouse(enableMouse), nil
}

package tviewyaml

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/cassdeckard/tviewyaml/builder"
	"github.com/cassdeckard/tviewyaml/config"
	"github.com/cassdeckard/tviewyaml/template"
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
	rootConfig, err := loader.LoadRoot("root.yaml")
	if err != nil {
		return nil, err
	}

	// Validate root config
	validator := config.NewValidator()
	if err := validator.ValidateRoot(rootConfig); err != nil {
		return nil, err
	}

	// Create builder
	uiBuilder := builder.NewBuilder(ctx)

	// Build all pages from config
	for _, pageRef := range rootConfig.Root.Pages {
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

	// Global keyboard shortcuts
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("main")
			return nil
		}
		return event
	})

	return app.SetRoot(pages, true).EnableMouse(true), nil
}

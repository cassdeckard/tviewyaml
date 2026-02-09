package tviewyaml

import (
	"fmt"
	"log"
	"time"

	"github.com/cassdeckard/tviewyaml/builder"
	"github.com/cassdeckard/tviewyaml/config"
	"github.com/cassdeckard/tviewyaml/template"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// AppBuilder provides a fluent API for building tview applications from YAML configuration
type AppBuilder struct {
	configDir string
	registry  *template.FunctionRegistry
}

// NewAppBuilder creates a new application builder
func NewAppBuilder(configDir string) *AppBuilder {
	return &AppBuilder{
		configDir: configDir,
		registry:  template.NewFunctionRegistry(),
	}
}

// WithTemplateFunction registers a custom template function
func (b *AppBuilder) WithTemplateFunction(name string, minArgs int, maxArgs *int, validator func(*template.Context, []string) error, handler interface{}) *AppBuilder {
	if err := b.registry.Register(name, minArgs, maxArgs, validator, handler); err != nil {
		// Log the error but continue building (could also panic or store errors)
		log.Printf("Warning: failed to register template function %q: %v", name, err)
	}
	return b
}

// RegisterTemplateFunctions calls fn with the builder so the app can register custom
// template functions (e.g. clock). Returns fn(b) for chaining.
func (b *AppBuilder) RegisterTemplateFunctions(fn func(*AppBuilder) *AppBuilder) *AppBuilder {
	return fn(b)
}

// Build creates and configures a tview application from YAML configuration files.
// Returns (app, pageErrors, err) where err is fatal (app config load/validate failure),
// and pageErrors are non-fatal per-page failures (missing/invalid pages are skipped).
func (b *AppBuilder) Build() (*tview.Application, []error, error) {
	// Initialize tview application
	app := tview.NewApplication()
	pages := tview.NewPages()

	// Create template context
	ctx := template.NewContext(app, pages)

	// Load configuration
	loader := config.NewLoader(b.configDir)
	appConfig, err := loader.LoadApp("app.yaml")
	if err != nil {
		return nil, nil, err
	}

	// Validate app config
	validator := config.NewValidator()
	if err := validator.ValidateApp(appConfig); err != nil {
		return nil, nil, err
	}
	if err := validator.ValidateAppRefs(appConfig, loader); err != nil {
		return nil, nil, err
	}

	// Create builder with registry
	uiBuilder := builder.NewBuilder(ctx, b.registry)

	// Build all pages from config, collecting non-fatal errors
	var pageErrors []error
	for _, pageRef := range appConfig.Application.Root.Pages {
		pageConfig, err := loader.LoadPage(pageRef.Ref)
		if err != nil {
			pageErrors = append(pageErrors, fmt.Errorf("error loading page %s: %w", pageRef.Name, err))
			continue
		}

		// Validate page config
		if err := validator.ValidatePage(pageConfig); err != nil {
			pageErrors = append(pageErrors, fmt.Errorf("invalid page config %s: %w", pageRef.Name, err))
			continue
		}

		pagePrimitive, err := uiBuilder.BuildFromConfig(pageConfig)
		if err != nil {
			pageErrors = append(pageErrors, fmt.Errorf("error building page %s: %w", pageRef.Name, err))
			continue
		}

		// Add to pages
		visible := pageRef.Name == "main"
		pages.AddPage(pageRef.Name, pagePrimitive, true, visible)
	}

	// Background goroutine: periodically refresh bound views whose state is dirty.
	// Does not depend on clock or user input; runs continuously and queues updates via QueueUpdateDraw.
	go func() {
		ticker := time.NewTicker(150 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			if !ctx.HasDirtyKeys() {
				continue
			}
			app.QueueUpdateDraw(func() {
				ctx.RefreshDirtyBoundViews()
			})
		}
	}()

	// Set input capture only when we have global key bindings; avoid running refresh
	// from capture to prevent deadlock (QueueUpdate would block) or draw re-entrancy.
	executor := template.NewExecutor(ctx, b.registry)
	ctx.SetExecutor(executor)
	if len(appConfig.Application.GlobalKeyBindings) > 0 {
		passthrough := appConfig.Application.EscapePassthroughPages
		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			// On Escape, if current page is in passthrough list, let the primitive (e.g. form) handle it.
			if event.Key() == tcell.KeyEscape && len(passthrough) > 0 {
				if front, _ := pages.GetFrontPage(); front != "" {
					for _, p := range passthrough {
						if p == front {
							return event
						}
					}
				}
			}
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

	// Apply mouse setting (default to true when not specified in config)
	enableMouse := true
	if appConfig.Application.EnableMouse != nil {
		enableMouse = *appConfig.Application.EnableMouse
	}

	return app.SetRoot(pages, true).EnableMouse(enableMouse), pageErrors, nil
}

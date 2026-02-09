package tviewyaml

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cassdeckard/tviewyaml/builder"
	"github.com/cassdeckard/tviewyaml/config"
	"github.com/cassdeckard/tviewyaml/template"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Application wraps tview.Application with lifecycle management for background goroutines
type Application struct {
	*tview.Application
	stopRefresh chan struct{}
}

// Stop gracefully shuts down the application and stops all background goroutines
func (a *Application) Stop() {
	if a.stopRefresh != nil {
		close(a.stopRefresh)
	}
	if a.Application != nil {
		a.Application.Stop()
	}
}

// AppBuilder provides a fluent API for building tview applications from YAML configuration
type AppBuilder struct {
	configDir string
	registry  *template.FunctionRegistry
	errors    []error
}

// NewAppBuilder creates a new application builder
func NewAppBuilder(configDir string) *AppBuilder {
	return &AppBuilder{
		configDir: configDir,
		registry:  template.NewFunctionRegistry(),
		errors:    make([]error, 0),
	}
}

// WithTemplateFunction registers a custom template function
func (b *AppBuilder) WithTemplateFunction(name string, minArgs int, maxArgs *int, validator func(*template.Context, []string) error, handler interface{}) *AppBuilder {
	if err := b.registry.Register(name, minArgs, maxArgs, validator, handler); err != nil {
		b.errors = append(b.errors, fmt.Errorf("failed to register template function %q: %w", name, err))
	}
	return b
}

// With calls fn with the builder so the app can perform custom
// registration with the AppBuilder. Returns fn(b) for chaining.
func (b *AppBuilder) With(fn func(*AppBuilder) *AppBuilder) *AppBuilder {
	return fn(b)
}

// Build creates and configures a tview application from YAML configuration files.
// Returns (app, pageErrors, err) where err is fatal (app config load/validate failure),
// and pageErrors are non-fatal per-page failures (missing/invalid pages are skipped).
func (b *AppBuilder) Build() (*Application, []error, error) {
	// Check for builder configuration errors first
	if len(b.errors) > 0 {
		return nil, nil, fmt.Errorf("builder configuration errors: %v", b.errors)
	}

	// Initialize tview application
	tvApp := tview.NewApplication()
	pages := tview.NewPages()

	// Create template context
	ctx := template.NewContext(tvApp, pages)

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

	// Validate template expressions before building pages
	if err := b.validateTemplateExpressions(appConfig, loader); err != nil {
		return nil, nil, fmt.Errorf("template validation failed: %w", err)
	}

	// Create builder with registry
	uiBuilder := builder.NewBuilder(ctx, b.registry)
	uiBuilder.SetLoader(loader) // Enable nested pages support

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

	// Create wrapped application with lifecycle management
	stopRefresh := make(chan struct{})
	app := &Application{
		Application: tvApp,
		stopRefresh: stopRefresh,
	}

	// Background goroutine: periodically refresh bound views whose state is dirty.
	// Does not depend on clock or user input; runs continuously and queues updates via QueueUpdateDraw.
	// The goroutine stops when stopRefresh channel is closed (via app.Stop()).
	go func() {
		ticker := time.NewTicker(150 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopRefresh:
				return
			case <-ticker.C:
				if !ctx.HasDirtyKeys() {
					continue
				}
				tvApp.QueueUpdateDraw(func() {
					ctx.RefreshDirtyBoundViews()
				})
			}
		}
	}()

	// Set input capture only when we have global key bindings; avoid running refresh
	// from capture to prevent deadlock (QueueUpdate would block) or draw re-entrancy.
	executor := template.NewExecutor(ctx, b.registry)
	ctx.SetExecutor(executor)
	if len(appConfig.Application.GlobalKeyBindings) > 0 {
		passthrough := appConfig.Application.EscapePassthroughPages
		tvApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

	app.Application = tvApp.SetRoot(pages, true).EnableMouse(enableMouse)
	return app, pageErrors, nil
}

// validateTemplateExpressions validates that all template expressions reference existing functions/evaluators
func (b *AppBuilder) validateTemplateExpressions(appConfig *config.AppConfig, loader *config.Loader) error {
	var errors []string

	// Validate global key bindings
	for _, binding := range appConfig.Application.GlobalKeyBindings {
		if binding.Action != "" {
			if errs := b.validateExpression(binding.Action, "global key binding"); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		}
	}

	// Validate all pages
	for _, pageRef := range appConfig.Application.Root.Pages {
		pageConfig, err := loader.LoadPage(pageRef.Ref)
		if err != nil {
			// Page loading errors are already handled in Build(), skip here
			continue
		}

		// Validate page-level expressions
		if errs := b.validatePageExpressions(pageConfig, pageRef.Name); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("found %d template validation error(s):\n  - %s", len(errors), strings.Join(errors, "\n  - "))
	}
	return nil
}

// validateExpression validates a single template expression
func (b *AppBuilder) validateExpression(expr, context string) []string {
	if expr == "" {
		return nil
	}

	var errors []string
	
	// Extract template expressions (handles both {{ }} and bare expressions)
	expr = strings.TrimSpace(expr)
	expr = strings.TrimPrefix(expr, "{{")
	expr = strings.TrimSuffix(expr, "}}")
	expr = strings.TrimSpace(expr)

	if expr == "" {
		return nil
	}

	// Parse function name from expression
	re := regexp.MustCompile(`^(\w+)`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) < 2 {
		return errors
	}

	funcName := matches[1]

	// Check if it exists as either a function or evaluator
	if _, ok := b.registry.Get(funcName); !ok {
		if _, ok := b.registry.GetEvaluator(funcName); !ok {
			errors = append(errors, fmt.Sprintf("%s: unknown function/evaluator %q in expression %q", context, funcName, expr))
		}
	}

	return errors
}

// validatePageExpressions validates all template expressions in a page config
func (b *AppBuilder) validatePageExpressions(page *config.PageConfig, pageName string) []string {
	var errors []string

	// Validate page-level callbacks (only those that exist on PageConfig)
	if page.OnSubmit != "" {
		errors = append(errors, b.validateExpression(page.OnSubmit, fmt.Sprintf("page %q OnSubmit", pageName))...)
	}
	if page.OnCancel != "" {
		errors = append(errors, b.validateExpression(page.OnCancel, fmt.Sprintf("page %q OnCancel", pageName))...)
	}
	if page.OnNodeSelected != "" {
		errors = append(errors, b.validateExpression(page.OnNodeSelected, fmt.Sprintf("page %q OnNodeSelected", pageName))...)
	}

	// Validate list items
	for i, item := range page.ListItems {
		if item.OnSelected != "" {
			errors = append(errors, b.validateExpression(item.OnSelected, fmt.Sprintf("page %q listItem[%d]", pageName, i))...)
		}
	}

	// Validate form items
	for i, item := range page.FormItems {
		if item.OnSelected != "" {
			errors = append(errors, b.validateExpression(item.OnSelected, fmt.Sprintf("page %q formItem[%d]", pageName, i))...)
		}
		if item.OnChanged != "" {
			errors = append(errors, b.validateExpression(item.OnChanged, fmt.Sprintf("page %q formItem[%d]", pageName, i))...)
		}
	}

	// Validate nested primitives recursively
	for i, flexItem := range page.Items {
		if flexItem.Primitive != nil {
			errors = append(errors, b.validatePrimitiveExpressions(flexItem.Primitive, fmt.Sprintf("page %q item[%d]", pageName, i))...)
		}
	}

	return errors
}

// validatePrimitiveExpressions validates all template expressions in a primitive (recursive)
func (b *AppBuilder) validatePrimitiveExpressions(prim *config.Primitive, context string) []string {
	var errors []string

	// Validate primitive callbacks
	if prim.OnSelected != "" {
		errors = append(errors, b.validateExpression(prim.OnSelected, fmt.Sprintf("%s OnSelected", context))...)
	}
	if prim.OnSubmit != "" {
		errors = append(errors, b.validateExpression(prim.OnSubmit, fmt.Sprintf("%s OnSubmit", context))...)
	}
	if prim.OnCancel != "" {
		errors = append(errors, b.validateExpression(prim.OnCancel, fmt.Sprintf("%s OnCancel", context))...)
	}
	if prim.OnChanged != "" {
		errors = append(errors, b.validateExpression(prim.OnChanged, fmt.Sprintf("%s OnChanged", context))...)
	}
	if prim.OnCellSelected != "" {
		errors = append(errors, b.validateExpression(prim.OnCellSelected, fmt.Sprintf("%s OnCellSelected", context))...)
	}
	if prim.OnNodeSelected != "" {
		errors = append(errors, b.validateExpression(prim.OnNodeSelected, fmt.Sprintf("%s OnNodeSelected", context))...)
	}

	// Validate list items
	for i, item := range prim.ListItems {
		if item.OnSelected != "" {
			errors = append(errors, b.validateExpression(item.OnSelected, fmt.Sprintf("%s listItem[%d]", context, i))...)
		}
	}

	// Validate form items
	for i, item := range prim.FormItems {
		if item.OnSelected != "" {
			errors = append(errors, b.validateExpression(item.OnSelected, fmt.Sprintf("%s formItem[%d]", context, i))...)
		}
		if item.OnChanged != "" {
			errors = append(errors, b.validateExpression(item.OnChanged, fmt.Sprintf("%s formItem[%d]", context, i))...)
		}
	}

	// Recurse into nested primitives
	for i, flexItem := range prim.Items {
		if flexItem.Primitive != nil {
			errors = append(errors, b.validatePrimitiveExpressions(flexItem.Primitive, fmt.Sprintf("%s flexItem[%d]", context, i))...)
		}
	}

	for i, gridItem := range prim.GridItems {
		if gridItem.Primitive != nil {
			errors = append(errors, b.validatePrimitiveExpressions(gridItem.Primitive, fmt.Sprintf("%s gridItem[%d]", context, i))...)
		}
	}

	return errors
}

package template

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rivo/tview"
)

// registerBuiltinFunctions registers all built-in template functions
func registerBuiltinFunctions(registry *FunctionRegistry) {
	// Helper to convert int to *int for maxArgs
	intPtr := func(i int) *int { return &i }

	// bindState: evaluator that returns current state value as string
	registry.RegisterEvaluator("bindState", 1, 1, func(ctx *Context, args []string) string {
		v, ok := ctx.GetState(args[0])
		if !ok {
			return ""
		}
		return fmt.Sprint(v)
	})

	// showNotification: sets notification state so bound TextViews display it.
	// Uses SetStateDirect (not SetState) because it's called from event handlers.
	registry.Register("showNotification", 1, intPtr(1), nil, func(ctx *Context, msg string) {
		ctx.SetStateDirect("notification", msg)
	})

	// switchToPage: switches to a different page
	registry.Register("switchToPage", 1, intPtr(1), nil, func(ctx *Context, pageName string) {
		ctx.Pages.SwitchToPage(pageName)
	})

	// removePage: removes a page from the pages container
	registry.Register("removePage", 1, intPtr(1), nil, func(ctx *Context, pageName string) {
		ctx.Pages.RemovePage(pageName)
	})

	// stopApp: stops the tview application
	registry.Register("stopApp", 0, intPtr(0), nil, func(ctx *Context) {
		ctx.App.Stop()
	})

	// showSimpleModal: displays a simple modal with text and buttons
	// First arg is text, remaining args are button labels (defaults to "OK" if no buttons provided)
	// Uses a unique page name so multiple modals can be shown without overwriting.
	registry.Register("showSimpleModal", 1, nil, nil, func(ctx *Context, args []string) {
		text := args[0]
		buttons := args[1:]
		if len(buttons) == 0 {
			buttons = []string{"OK"}
		}

		pageName := "simple-modal-" + strconv.FormatInt(time.Now().UnixNano(), 10)
		modal := tview.NewModal().
			SetText(text).
			AddButtons(buttons).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				ctx.Pages.RemovePage(pageName)
			})
		ctx.Pages.AddPage(pageName, modal, false, true)
	})

	// noop: does nothing (useful for testing or placeholder actions)
	registry.Register("noop", 0, intPtr(0), nil, func(ctx *Context) {
		// Do nothing
	})
}

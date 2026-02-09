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

	// showSimpleModal: displays a simple modal with text and buttons.
	// Args: text, [button labels...], [optional onDone template]. Example: "Done!" "OK" "switchToPage \"main\""
	// Uses a unique page name so multiple modals can be shown without overwriting.
	registry.Register("showSimpleModal", 1, nil, nil, func(ctx *Context, args []string) {
		text := args[0]
		var buttons []string
		var onDone string
		if len(args) >= 3 {
			onDone = args[len(args)-1]
			buttons = args[1 : len(args)-1]
		} else {
			buttons = args[1:]
		}
		if len(buttons) == 0 {
			buttons = []string{"OK"}
		}

		pageName := "simple-modal-" + strconv.FormatInt(time.Now().UnixNano(), 10)
		modal := tview.NewModal().
			SetText(text).
			AddButtons(buttons).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if onDone != "" {
					ctx.RunCallback(onDone)
				}
				ctx.Pages.RemovePage(pageName)
			})
		ctx.Pages.AddPage(pageName, modal, false, true)
	})

	// runFormSubmit: runs the submit callback registered for the form name (e.g. from a Submit button).
	registry.Register("runFormSubmit", 1, intPtr(1), nil, func(ctx *Context, formName string) {
		ctx.RunFormSubmit(formName)
	})

	// showSelectedCellModal: shows a modal with the currently selected table cell info (reads __selectedCellText, __selectedRow, __selectedCol from state).
	registry.Register("showSelectedCellModal", 0, intPtr(0), nil, func(ctx *Context) {
		cellText, _ := ctx.GetState("__selectedCellText")
		row, _ := ctx.GetState("__selectedRow")
		col, _ := ctx.GetState("__selectedCol")
		text := fmt.Sprintf("Selected cell: %s\nRow: %v, Column: %v", cellText, row, col)
		pageName := "cell-modal-" + strconv.FormatInt(time.Now().UnixNano(), 10)
		modal := tview.NewModal().
			SetText(text).
			AddButtons([]string{"OK"}).
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

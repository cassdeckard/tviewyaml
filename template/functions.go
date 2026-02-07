package template

import (
	"github.com/rivo/tview"
)

// CreateFuncMap creates a whitelist of safe template functions
func CreateFuncMap(ctx *Context) map[string]interface{} {
	return map[string]interface{}{
		// Navigation functions
		"switchToPage": func(name string) func() {
			return func() { ctx.Pages.SwitchToPage(name) }
		},
		"removePage": func(name string) func() {
			return func() { ctx.Pages.RemovePage(name) }
		},
		"stopApp": func() func() {
			return func() { ctx.App.Stop() }
		},

		// Simple modal with customizable buttons
		// First parameter is the text to display
		// Remaining parameters are button labels
		"showSimpleModal": func(text string, buttons ...string) func() {
			return func() {
				// If no buttons provided, default to "OK"
				if len(buttons) == 0 {
					buttons = []string{"OK"}
				}
				
				modal := tview.NewModal().
					SetText(text).
					AddButtons(buttons).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						// Just remove the modal overlay, the underlying page should still be visible
						ctx.Pages.RemovePage("simple-modal")
					})
				// Add as overlay (false = not fullscreen, true = visible)
				ctx.Pages.AddPage("simple-modal", modal, false, true)
			}
		},

		// Color helper
		"color": ctx.Colors.Parse,

		// Utility functions
		"noop": func() func() { return func() {} },
	}
}

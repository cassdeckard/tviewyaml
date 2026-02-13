package app

import (
	"fmt"

	"github.com/cassdeckard/tviewyaml"
	"github.com/cassdeckard/tviewyaml/template"
	"github.com/rivo/tview"
)

// RegisterDynamicPages adds dynamic page creation functions as custom template functions
func RegisterDynamicPages(b *tviewyaml.AppBuilder) *tviewyaml.AppBuilder {
	maxOne := 1
	b.WithTemplateFunction("createDetailPageFor", 1, &maxOne, nil, createDetailPageFor)
	return b
}

// createDetailPageFor creates a page dynamically for the given item name
func createDetailPageFor(ctx *template.Context, itemID string) {

	// Build a page manually using tview
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true)
	flex.SetTitle(fmt.Sprintf("Dynamic Detail: %s", itemID))

	// Add header explaining what happened
	header := tview.NewTextView()
	header.SetText(fmt.Sprintf(
		"[yellow::b]This page was created dynamically at runtime![::-]\n\n"+
			"[green]Item:[::-] %s\n"+
			"[green]Page Name:[::-] dynamic-detail\n"+
			"[green]Created:[::-] Programmatically in Go code\n\n"+
			"[cyan]How it works:[::-]\n"+
			"1. Button calls: createDetailPageFor \"%s\"\n"+
			"2. Go function receives item name as parameter\n"+
			"3. Function builds tview.Flex with content\n"+
			"4. Adds to Pages container via ctx.Pages.AddPage\n"+
			"5. Switches to the new page with SwitchToPage",
		itemID, itemID,
	))
	header.SetDynamicColors(true)
	header.SetTextAlign(tview.AlignLeft)
	flex.AddItem(header, 0, 1, false)

	// Add action buttons
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Actions")

	form.AddButton("Remove This Page & Return", func() {
		if ctx.Pages != nil {
			ctx.Pages.RemovePage("dynamic-detail")
			ctx.Pages.SwitchToPage("dynamic-pages")
		}
	})

	form.AddButton("Create Another Page", func() {
		// Switch back to dynamic-pages to select another item
		if ctx.Pages != nil {
			ctx.Pages.SwitchToPage("dynamic-pages")
		}
	})

	form.AddButton("Main Menu", func() {
		if ctx.Pages != nil {
			ctx.Pages.RemovePage("dynamic-detail")
			ctx.Pages.SwitchToPage("main")
		}
	})

	flex.AddItem(form, 9, 0, true)

	// Add to pages container and switch to it
	if ctx.Pages != nil {
		// Check if page already exists and remove it first
		ctx.Pages.RemovePage("dynamic-detail")
		ctx.Pages.AddPage("dynamic-detail", flex, true, false)
		ctx.Pages.SwitchToPage("dynamic-detail")
	}
}

package builder

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/cassdeckard/tviewyaml/config"
)

// Factory creates tview primitives based on configuration
type Factory struct{}

// NewFactory creates a new primitive factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreatePrimitive creates a tview primitive based on type
func (f *Factory) CreatePrimitive(prim *config.Primitive) (tview.Primitive, error) {
	switch prim.Type {
	case "box":
		return tview.NewBox(), nil

	case "textView":
		return tview.NewTextView(), nil

	case "button":
		label := prim.Label
		if label == "" {
			label = "Button"
		}
		return tview.NewButton(label), nil

	case "list":
		return tview.NewList(), nil

	case "flex":
		flex := tview.NewFlex()
		if prim.Direction == "row" {
			flex.SetDirection(tview.FlexRow)
		}
		return flex, nil

	case "form":
		return tview.NewForm(), nil

	case "inputField":
		return tview.NewInputField(), nil

	case "checkbox":
		return tview.NewCheckbox(), nil

	case "dropdown":
		return tview.NewDropDown(), nil

	case "table":
		return tview.NewTable(), nil

	case "textArea":
		return tview.NewTextArea(), nil

	case "modal":
		return tview.NewModal(), nil

	case "pages":
		return tview.NewPages(), nil

	case "grid":
		return tview.NewGrid(), nil

	case "treeView":
		return tview.NewTreeView(), nil

	default:
		return nil, fmt.Errorf("unknown primitive type: %s", prim.Type)
	}
}

// CreatePrimitiveFromPageConfig creates a top-level primitive from a page config
func (f *Factory) CreatePrimitiveFromPageConfig(cfg *config.PageConfig) (tview.Primitive, error) {
	switch cfg.Type {
	case "list":
		return tview.NewList(), nil

	case "flex":
		flex := tview.NewFlex()
		if cfg.Direction == "row" {
			flex.SetDirection(tview.FlexRow)
		}
		return flex, nil

	case "form":
		return tview.NewForm(), nil

	case "table":
		return tview.NewTable(), nil

	case "grid":
		return tview.NewGrid(), nil

	case "pages":
		return tview.NewPages(), nil

	default:
		return nil, fmt.Errorf("unknown page type: %s", cfg.Type)
	}
}

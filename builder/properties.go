package builder

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/cassdeckard/tviewyaml/config"
	"github.com/cassdeckard/tviewyaml/template"
)

// PropertyMapper applies YAML properties to tview primitives
type PropertyMapper struct {
	colorHelper *template.ColorHelper
	context     *template.Context
	executor    *template.Executor
}

// NewPropertyMapper creates a new property mapper
func NewPropertyMapper(ctx *template.Context, executor *template.Executor) *PropertyMapper {
	return &PropertyMapper{
		colorHelper: &template.ColorHelper{},
		context:     ctx,
		executor:    executor,
	}
}

// ApplyProperties applies configuration properties to a primitive
func (pm *PropertyMapper) ApplyProperties(primitive tview.Primitive, prim *config.Primitive) error {
	// Common properties that apply to Box (base of most primitives)
	if b, ok := primitive.(interface {
		SetBorder(bool) *tview.Box
		SetTitle(string) *tview.Box
		SetTitleAlign(int) *tview.Box
	}); ok {
		if prim.Border {
			b.SetBorder(true)
		}
		if prim.Title != "" {
			b.SetTitle(prim.Title)
		}
		if prim.TitleAlign != "" {
			b.SetTitleAlign(template.ParseAlignment(prim.TitleAlign))
		}
	}

	// Type-specific properties
	switch v := primitive.(type) {
	case *tview.TextView:
		return pm.applyTextViewProperties(v, prim)
	case *tview.Button:
		return pm.applyButtonProperties(v, prim)
	case *tview.InputField:
		return pm.applyInputFieldProperties(v, prim)
	case *tview.Checkbox:
		return pm.applyCheckboxProperties(v, prim)
	case *tview.DropDown:
		return pm.applyDropDownProperties(v, prim)
	case *tview.Table:
		return pm.applyTableProperties(v, prim)
	case *tview.List:
		return pm.applyListProperties(v, prim)
	}

	return nil
}

func (pm *PropertyMapper) applyTextViewProperties(tv *tview.TextView, prim *config.Primitive) error {
	if prim.Text != "" {
		if strings.Contains(prim.Text, "{{") && strings.Contains(prim.Text, "}}") && pm.executor != nil {
			// Template syntax: evaluate once and register for deferred refresh on key events
			result, err := pm.executor.EvaluateToString(prim.Text)
			if err != nil {
				return fmt.Errorf("template evaluation failed: %w", err)
			}
			tv.SetText(result)
			keys := pm.executor.ExtractBindStateKeys(prim.Text)
			templateStr := prim.Text
			for _, key := range keys {
				pm.context.RegisterBoundView(key, template.BoundView{
					Refresh: func() string {
						s, err := pm.executor.EvaluateToString(templateStr)
						if err != nil {
							return ""
						}
						return s
					},
					SetText: func(s string) { tv.SetText(s) },
				})
			}
		} else {
			tv.SetText(prim.Text)
		}
	}
	if prim.TextAlign != "" {
		tv.SetTextAlign(template.ParseAlignment(prim.TextAlign))
	}
	if prim.TextColor != "" {
		tv.SetTextColor(pm.colorHelper.Parse(prim.TextColor))
	}
	// Enable dynamic colors and regions if specified
	if prim.DynamicColors {
		tv.SetDynamicColors(true)
	}
	if prim.Regions {
		tv.SetRegions(true)
		
		// Add region navigation handlers
		tv.SetDoneFunc(func(key tcell.Key) {
			currentSelection := tv.GetHighlights()
			if key == tcell.KeyEnter {
				if len(currentSelection) > 0 {
					tv.Highlight()
				} else {
					tv.Highlight("0").ScrollToHighlight()
				}
			} else if len(currentSelection) > 0 {
				index := 0
				fmt.Sscanf(currentSelection[0], "%d", &index)
				if key == tcell.KeyTab {
					index = (index + 1) % 3
				} else if key == tcell.KeyBacktab {
					index = (index - 1 + 3) % 3
				} else {
					return
				}
				tv.Highlight(fmt.Sprintf("%d", index)).ScrollToHighlight()
			}
		})
	}
	return nil
}

func (pm *PropertyMapper) applyButtonProperties(btn *tview.Button, prim *config.Primitive) error {
	// Button label is set in factory
	return nil
}

func (pm *PropertyMapper) applyInputFieldProperties(input *tview.InputField, prim *config.Primitive) error {
	if prim.Label != "" {
		input.SetLabel(prim.Label)
	}
	if prim.Text != "" {
		input.SetText(prim.Text)
	}
	return nil
}

func (pm *PropertyMapper) applyCheckboxProperties(cb *tview.Checkbox, prim *config.Primitive) error {
	if prim.Label != "" {
		cb.SetLabel(prim.Label)
	}
	if prim.Checked {
		cb.SetChecked(true)
	}
	// Add changed handler to trigger redraws (needed for standalone checkboxes)
	cb.SetChangedFunc(func(checked bool) {
		if pm.context != nil && pm.context.App != nil {
			pm.context.App.Draw()
		}
	})
	return nil
}

func (pm *PropertyMapper) applyDropDownProperties(dd *tview.DropDown, prim *config.Primitive) error {
	if prim.Label != "" {
		dd.SetLabel(prim.Label)
	}
	if len(prim.Options) > 0 {
		dd.SetOptions(prim.Options, nil)
	}
	return nil
}

func (pm *PropertyMapper) applyTableProperties(table *tview.Table, prim *config.Primitive) error {
	// Table borders, fixed rows/columns, and data are all handled in populateTableData
	return nil
}

func (pm *PropertyMapper) applyListProperties(list *tview.List, prim *config.Primitive) error {
	// List items are handled in builder
	return nil
}

// ApplyPageProperties applies page-level properties to a primitive
func (pm *PropertyMapper) ApplyPageProperties(primitive tview.Primitive, cfg *config.PageConfig) error {
	// Common properties
	if b, ok := primitive.(interface {
		SetBorder(bool) *tview.Box
		SetTitle(string) *tview.Box
		SetTitleAlign(int) *tview.Box
	}); ok {
		if cfg.Border {
			b.SetBorder(true)
		}
		if cfg.Title != "" {
			b.SetTitle(cfg.Title)
		}
		if cfg.TitleAlign != "" {
			b.SetTitleAlign(template.ParseAlignment(cfg.TitleAlign))
		}
	}

	// Type-specific properties
	switch v := primitive.(type) {
	case *tview.Flex:
		if cfg.Direction == "row" {
			v.SetDirection(tview.FlexRow)
		}
	}

	return nil
}

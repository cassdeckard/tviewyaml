package builder

import (
	"testing"

	"github.com/cassdeckard/tviewyaml/config"
	"github.com/cassdeckard/tviewyaml/template"
	"github.com/rivo/tview"
)

func TestApplyTextViewProperties_OnHighlighted(t *testing.T) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	ctx := template.NewContext(app, pages)
	registry := template.NewFunctionRegistry()
	executor := template.NewExecutor(ctx, registry)
	pm := NewPropertyMapper(ctx, executor)

	tv := tview.NewTextView()
	tv.SetRegions(true)
	tv.SetText(`["a"]Region A[""]`)

	prim := &config.Primitive{
		Type:          "textView",
		Regions:       true,
		OnHighlighted: `{{ showNotification "highlighted" }}`,
	}

	err := pm.ApplyProperties(tv, prim)
	if err != nil {
		t.Fatalf("ApplyProperties: %v", err)
	}
	// Wiring succeeded; SetHighlightedFunc is attached. The callback runs when
	// tview draws after a highlight change (tested via acceptance if needed).
}

func TestBuildFlex_TextViewWithOnHighlighted(t *testing.T) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	ctx := template.NewContext(app, pages)
	registry := template.NewFunctionRegistry()
	b := NewBuilder(ctx, registry)

	pageConfig := &config.PageConfig{
		Type: "flex",
		Items: []config.FlexItem{
			{
				Primitive: &config.Primitive{
					Type:          "textView",
					Regions:       true,
					DynamicColors: true,
					Text:          `["slide1"]Slide 1[""] ["slide2"]Slide 2[""]`,
					OnHighlighted: `{{ switchToPage "main" }}`,
				},
				FixedSize:  0,
				Proportion: 1,
				Focus:     true,
			},
		},
	}

	result, err := b.BuildFromConfig(pageConfig)
	if err != nil {
		t.Fatalf("BuildFromConfig: %v", err)
	}

	flex, ok := result.(*tview.Flex)
	if !ok {
		t.Fatalf("expected *tview.Flex, got %T", result)
	}
	if flex.GetItemCount() != 1 {
		t.Errorf("GetItemCount() = %d, want 1", flex.GetItemCount())
	}
}

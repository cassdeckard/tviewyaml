package builder

import (
	"testing"

	"github.com/cassdeckard/tviewyaml/config"
	"github.com/cassdeckard/tviewyaml/template"
	"github.com/rivo/tview"
)

func TestBuildFlex_SpacerItems(t *testing.T) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	ctx := template.NewContext(app, pages)
	registry := template.NewFunctionRegistry()
	b := NewBuilder(ctx, registry)

	pageConfig := &config.PageConfig{
		Type:      "flex",
		Direction: "row",
		Items: []config.FlexItem{
			{
				Primitive:  &config.Primitive{Type: "textView", Text: "Left"},
				FixedSize:  10,
				Proportion: 0,
				Focus:      false,
			},
			{
				Primitive:  nil, // spacer
				FixedSize:  0,
				Proportion: 1,
				Focus:      false,
			},
			{
				Primitive:  &config.Primitive{Type: "textView", Text: "Right"},
				FixedSize:  15,
				Proportion: 0,
				Focus:      true,
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

	if got := flex.GetItemCount(); got != 3 {
		t.Errorf("GetItemCount() = %d, want 3", got)
	}

	// Middle item should be nil (spacer)
	if got := flex.GetItem(1); got != nil {
		t.Errorf("GetItem(1) = %v, want nil (spacer)", got)
	}
}

func TestBuildFlex_SpacerFlag(t *testing.T) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	ctx := template.NewContext(app, pages)
	registry := template.NewFunctionRegistry()
	b := NewBuilder(ctx, registry)

	pageConfig := &config.PageConfig{
		Type:      "flex",
		Direction: "row",
		Items: []config.FlexItem{
			{
				Primitive:  &config.Primitive{Type: "textView", Text: "A"},
				FixedSize:  5,
				Proportion: 0,
				Focus:      false,
			},
			{
				Spacer:     true, // explicit spacer
				FixedSize:  0,
				Proportion: 1,
				Focus:      false,
			},
			{
				Primitive:  &config.Primitive{Type: "textView", Text: "B"},
				FixedSize:  5,
				Proportion: 0,
				Focus:      true,
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

	if got := flex.GetItemCount(); got != 3 {
		t.Errorf("GetItemCount() = %d, want 3", got)
	}

	if got := flex.GetItem(1); got != nil {
		t.Errorf("GetItem(1) = %v, want nil (spacer)", got)
	}
}

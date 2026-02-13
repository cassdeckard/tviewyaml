package acceptance_test

import (
	"testing"
)

func TestAcceptance_LayoutAtMultipleSizes(t *testing.T) {
	sizes := []struct {
		name       string
		cols, rows int
	}{
		{"80x24", 80, 24},
		{"120x30", 120, 30},
		{"40x10", 40, 10},
	}
	for _, sz := range sizes {
		t.Run(sz.name, func(t *testing.T) {
			h := newAcceptanceHarness(t, sz.cols, sz.rows)
			defer h.stop()
			// At 40 cols the full title is truncated; at 80+ "Tview Feature Demos" is visible.
			if !h.screenContains("Feature Demos") {
				t.Errorf("screen should contain main title (e.g. Feature Demos); content snippet: %s",
					truncate(h.getContent(), 500))
			}
			if !h.screenContains("Box") {
				t.Errorf("screen should contain %q", "Box")
			}
			if !h.screenContains("Button") {
				t.Errorf("screen should contain %q", "Button")
			}
			// Snapshot comparison (name derived from t.Name() -> e.g. TestAcceptance_LayoutAtMultipleSizes_80x24.terminal)
			h.AssertSnapshot(t, "")
		})
	}
}

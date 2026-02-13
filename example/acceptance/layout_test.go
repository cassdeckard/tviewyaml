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
			if !h.screenContains("Tview Feature Demos") {
				t.Errorf("screen should contain main title %q; content snippet: %s",
					"Tview Feature Demos", truncate(h.getContent(), 500))
			}
			if !h.screenContains("Box") {
				t.Errorf("screen should contain %q", "Box")
			}
			if !h.screenContains("Button") {
				t.Errorf("screen should contain %q", "Button")
			}
		})
	}
}

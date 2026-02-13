package acceptance_test

import (
	"testing"
)

func TestAcceptance_LayoutAtMultipleSizes(t *testing.T) {
	runAtSizes(t, func(t *testing.T, h *acceptanceHarness) {
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
		h.AssertSnapshot(t, "")
	})
}

package acceptance_test

import (
	"testing"
)

func TestAcceptance_KeyNavigation(t *testing.T) {
	runAtSizes(t, func(t *testing.T, h *acceptanceHarness) {
		t.Run("MainMenu", func(t *testing.T) {
			h.AssertSnapshot(t, "")
		})

		// Shortcut 'b' goes to Box page (from main.yaml). Wait for 2 draws so we see the new page.
		h.typeKey("b")
		if !h.waitForDraws(2) {
			t.Fatal("timeout waiting for draws after pressing b")
		}
		t.Run("BoxPage", func(t *testing.T) {
			if !h.screenContains("Box") {
				t.Errorf("after pressing b, screen should show Box page; content snippet: %s",
					truncate(h.getContent(), 500))
			}
			// Box page title from example/config/box.yaml
			if !h.screenContains("Box Demo") {
				t.Errorf("screen should contain Box page title; content snippet: %s",
					truncate(h.getContent(), 500))
			}
			h.AssertSnapshot(t, "")
		})

		// Escape returns to main (global keybinding).
		h.typeKey("Escape")
		if !h.waitForDraw() {
			t.Fatal("timeout waiting for draw after Escape")
		}
		t.Run("BackToMain", func(t *testing.T) {
			// At 40 cols the full title is truncated; "Feature Demos" is visible at all sizes.
			if !h.screenContains("Feature Demos") {
				t.Errorf("after Escape, screen should show main menu; content snippet: %s",
					truncate(h.getContent(), 500))
			}
			h.AssertSnapshot(t, "")
		})
	})
}

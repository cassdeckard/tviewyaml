package acceptance_test

import (
	"testing"
)

func TestAcceptance_KeyNavigation(t *testing.T) {
	h := newAcceptanceHarness(t, 80, 24)
	defer h.stop()

	// Shortcut 'b' goes to Box page (from main.yaml). Wait for 2 draws so we see the new page.
	h.typeKey("b")
	if !h.waitForDraws(2) {
		t.Fatal("timeout waiting for draws after pressing b")
	}
	if !h.screenContains("Box") {
		t.Errorf("after pressing b, screen should show Box page; content snippet: %s",
			truncate(h.getContent(), 500))
	}
	// Box page title from example/config/box.yaml
	if !h.screenContains("Box Demo") {
		t.Errorf("screen should contain Box page title; content snippet: %s",
			truncate(h.getContent(), 500))
	}

	// Escape returns to main (global keybinding).
	h.typeKey("Escape")
	if !h.waitForDraw() {
		t.Fatal("timeout waiting for draw after Escape")
	}
	if !h.screenContains("Tview Feature Demos") {
		t.Errorf("after Escape, screen should show main menu; content snippet: %s",
			truncate(h.getContent(), 500))
	}
}

package acceptance_test

import (
	"testing"
)

// TestAcceptance_Snapshot asserts the initial main menu against a golden snapshot.
// Create or update the golden file with: UPDATE_TERMINAL_SNAPSHOTS=1 go test ./example/acceptance/ -run TestAcceptance_Snapshot
func TestAcceptance_Snapshot(t *testing.T) {
	h := newAcceptanceHarness(t, 80, 24)
	defer h.stop()

	h.AssertSnapshot(t, "main_menu")
}

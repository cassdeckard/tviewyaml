package acceptance_test

import (
	"testing"
)

// navPages defines main-menu shortcuts and page-specific assertions (from example/config/main.yaml).
// escapeExtra: keys to press after Escape when page has modal (e.g. Form's onCancel shows "OK" modal).
// Use navKey for global shortcuts (e.g. "Alt+6") when list shortcut is unreliable at small sizes.
var navPages = []struct {
	key         string
	navKey      string // if set, use instead of key (for Alt+N global shortcuts)
	subtest     string
	contains    string // distinctive content to assert (visible at 40x10)
	escapeExtra string // extra keys after Escape (e.g. "Enter" for Form's cancel modal)
}{
	{"b", "", "BoxPage", "Box Demo", ""},
	{"u", "", "ButtonPage", "Button Demo", ""},
	{"c", "", "CheckboxPage", "Checkbox Demo", ""},
	{"i", "", "InputFieldPage", "InputField Demo", ""},
	{"f", "", "FormPage", "Form Demo", "Enter"}, // Escape triggers onCancel modal, Enter dismisses
	{"l", "", "ListPage", "List Demo", ""},
	{"t", "", "TablePage", "Table Demo", ""},
	{"v", "", "TextViewPage", "TextView Demo", ""},
	{"r", "", "TreeViewPage", "TreeView Demo", ""},
	{"e", "", "TreeViewStandalonePage", "TreeView as Page", ""},
	{"s", "", "TreeViewModesPage", "Parent (auto)", ""},
	{"d", "", "DropDownPage", "DropDown Demo", ""},
	{"m", "", "ModalPage", "YAML-Configured", ""},
	{"y", "", "DynamicPagesPage", "Dynamic Page", ""},
	{"n", "", "NestedPagesPage", "Nested Pages", ""},
	{"x", "", "FlexPage", "Flex Demo", ""},
	{"g", "", "GridPage", "Grid Demo", ""},
	{"k", "Alt+6", "ClockPage", "Time:", ""}, // Alt+6 more reliable; "Time:" is distinctive (state display)
	{"w", "End", "StateBindingPage", "Reactive", ""}, // End key (global); "Reactive" in "Reactive State Pattern" visible at 40x10
	{"h", "Meta+H", "HelpPage", "Alt+4", "Enter"}, // Meta+H (global); "Alt+4" in help content visible at all sizes
}

func TestAcceptance_KeyNavigation(t *testing.T) {
	runAtSizes(t, func(t *testing.T, h *acceptanceHarness) {
		t.Run("MainMenu", func(t *testing.T) {
			h.AssertSnapshot(t, "")
		})

		for _, p := range navPages {
			key := p.key
			if p.navKey != "" {
				key = p.navKey
			}
			h.typeKey(key)
			if !h.waitForContent(p.contains) {
				t.Fatalf("timeout waiting for %q after pressing %q; content snippet: %s",
					p.contains, key, truncate(h.getContent(), 500))
			}
			t.Run(p.subtest, func(t *testing.T) {
				if !h.screenContains(p.contains) {
					t.Errorf("after pressing %q, screen should contain %q; content snippet: %s",
						key, p.contains, truncate(h.getContent(), 500))
				}
				h.AssertSnapshot(t, "")
			})
			h.typeKey("Escape")
			if p.escapeExtra != "" {
				h.typeKey(p.escapeExtra)
				if !h.waitForDraw() {
					t.Fatalf("timeout waiting for draw after %q from %s", p.escapeExtra, p.subtest)
				}
			}
			if !h.waitForContent("Feature Demos") {
				t.Fatalf("timeout waiting for main menu after Escape from %s; content snippet: %s",
					p.subtest, truncate(h.getContent(), 500))
			}
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

package keys

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestParseKey(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantKey     tcell.Key
		wantMod     tcell.ModMask
		wantRune    rune
		wantErr     bool
		errContains string
	}{
		// Special keys
		{"Escape", "Escape", tcell.KeyEscape, 0, 0, false, ""},
		{"esc lowercase", "esc", tcell.KeyEscape, 0, 0, false, ""},
		{"Enter", "Enter", tcell.KeyEnter, 0, 0, false, ""},
		{"return", "return", tcell.KeyEnter, 0, 0, false, ""},
		{"Tab", "Tab", tcell.KeyTab, 0, 0, false, ""},
		{"Backtab", "Backtab", tcell.KeyBacktab, 0, 0, false, ""},
		{"Backspace", "Backspace", tcell.KeyBackspace, 0, 0, false, ""},
		{"bs", "bs", tcell.KeyBackspace, 0, 0, false, ""},
		{"Delete", "Delete", tcell.KeyDelete, 0, 0, false, ""},
		{"del", "del", tcell.KeyDelete, 0, 0, false, ""},
		{"Insert", "Insert", tcell.KeyInsert, 0, 0, false, ""},
		{"ins", "ins", tcell.KeyInsert, 0, 0, false, ""},
		{"Up", "Up", tcell.KeyUp, 0, 0, false, ""},
		{"Down", "Down", tcell.KeyDown, 0, 0, false, ""},
		{"Left", "Left", tcell.KeyLeft, 0, 0, false, ""},
		{"Right", "Right", tcell.KeyRight, 0, 0, false, ""},
		{"Home", "Home", tcell.KeyHome, 0, 0, false, ""},
		{"End", "End", tcell.KeyEnd, 0, 0, false, ""},
		{"PgUp", "PgUp", tcell.KeyPgUp, 0, 0, false, ""},
		{"PageUp", "PageUp", tcell.KeyPgUp, 0, 0, false, ""},
		{"PgDn", "PgDn", tcell.KeyPgDn, 0, 0, false, ""},
		{"PageDown", "PageDown", tcell.KeyPgDn, 0, 0, false, ""},

		// Function keys
		{"F1", "F1", tcell.KeyF1, 0, 0, false, ""},
		{"F1 lowercase", "f1", tcell.KeyF1, 0, 0, false, ""},
		{"F12", "F12", tcell.KeyF1+11, 0, 0, false, ""},
		{"F12 lowercase", "f12", tcell.KeyF1+11, 0, 0, false, ""},

		// Character keys
		{"single letter", "a", tcell.KeyRune, 0, 'a', false, ""},
		{"single letter uppercase", "Q", tcell.KeyRune, 0, 'Q', false, ""},
		{"space", "space", tcell.KeyRune, 0, ' ', false, ""},
		{"unicode rune", "ñ", tcell.KeyRune, 0, 'ñ', false, ""},

		// Modifiers
		{"Ctrl+Q", "Ctrl+Q", tcell.KeyRune, tcell.ModCtrl, 'Q', false, ""},
		{"Control+Q", "Control+Q", tcell.KeyRune, tcell.ModCtrl, 'Q', false, ""},
		{"Alt+F1", "Alt+F1", tcell.KeyF1, tcell.ModAlt, 0, false, ""},
		{"Shift+Tab", "Shift+Tab", tcell.KeyTab, tcell.ModShift, 0, false, ""},
		{"Meta+a", "Meta+a", tcell.KeyRune, tcell.ModMeta, 'a', false, ""},
		{"Ctrl+Alt+Q", "Ctrl+Alt+Q", tcell.KeyRune, tcell.ModCtrl|tcell.ModAlt, 'Q', false, ""},
		{"Ctrl+Shift+Enter", "Ctrl+Shift+Enter", tcell.KeyEnter, tcell.ModCtrl|tcell.ModShift, 0, false, ""},

		// Whitespace handling
		{"trimmed spaces", "  Escape  ", tcell.KeyEscape, 0, 0, false, ""},
		{"trimmed modifier spaces", " Ctrl + Q ", tcell.KeyRune, tcell.ModCtrl, 'Q', false, ""},

		// Error cases
		{"empty string", "", 0, 0, 0, true, "empty key string"},
		{"whitespace only", "   ", 0, 0, 0, true, "empty key string"},
		{"unknown modifier", "Invalid+Q", 0, 0, 0, true, "unknown modifier"},
		{"unknown key", "UnknownKey", 0, 0, 0, true, "unknown key"},
		{"F0 invalid", "F0", 0, 0, 0, true, "unknown key"},
		{"F13 invalid", "F13", 0, 0, 0, true, "unknown key"},
		{"F99 invalid", "F99", 0, 0, 0, true, "unknown key"},
		{"multiple runes", "ab", 0, 0, 0, true, "unknown key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotMod, gotRune, err := ParseKey(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseKey(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("ParseKey(%q) error = %v, want error containing %q", tt.input, err, tt.errContains)
					}
				}
				return
			}
			if gotKey != tt.wantKey {
				t.Errorf("ParseKey(%q) gotKey = %v, want %v", tt.input, gotKey, tt.wantKey)
			}
			if gotMod != tt.wantMod {
				t.Errorf("ParseKey(%q) gotMod = %v, want %v", tt.input, gotMod, tt.wantMod)
			}
			if gotRune != tt.wantRune {
				t.Errorf("ParseKey(%q) gotRune = %v, want %v", tt.input, gotRune, tt.wantRune)
			}
		})
	}
}

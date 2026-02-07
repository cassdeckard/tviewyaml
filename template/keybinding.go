package template

import (
	"unicode"

	"github.com/cassdeckard/tviewyaml/config"
	"github.com/cassdeckard/tviewyaml/keys"
	"github.com/gdamore/tcell/v2"
)

// MatchesKeyBinding returns true if the given key event matches the binding.
// The binding.Key string supports formats like "Escape", "Enter", "Ctrl+Q", "Ctrl+C", "F1", "Alt+a".
func MatchesKeyBinding(event *tcell.EventKey, binding config.KeyBinding) bool {
	key, mod, ch, err := keys.ParseKey(binding.Key)
	if err != nil {
		return false
	}
	if key == tcell.KeyRune && ch == 0 {
		return false
	}

	eventMod := event.Modifiers()
	eventKey := event.Key()
	eventRune := event.Rune()

	// Modifiers must match (only check the ones we care about - masked for resilience)
	wantMod := mod & (tcell.ModCtrl | tcell.ModAlt | tcell.ModShift | tcell.ModMeta)
	gotMod := eventMod & (tcell.ModCtrl | tcell.ModAlt | tcell.ModShift | tcell.ModMeta)
	if wantMod != gotMod {
		return false
	}

	if key == tcell.KeyRune {
		// Letter/key with optional modifier
		if eventKey != tcell.KeyRune {
			return false
		}
		// Ctrl+letter produces ASCII control codes (Ctrl+A=1, Ctrl+Q=17, etc.)
		if mod&tcell.ModCtrl != 0 && ch >= 'a' && ch <= 'z' {
			ctrlRune := rune(ch - 'a' + 1)
			return eventRune == ctrlRune || eventRune == ch
		}
		if mod&tcell.ModCtrl != 0 && ch >= 'A' && ch <= 'Z' {
			ctrlRune := rune(ch - 'A' + 1)
			return eventRune == ctrlRune || eventRune == ch || eventRune == unicode.ToLower(ch)
		}
		return eventRune == ch || eventRune == unicode.ToLower(ch) || eventRune == unicode.ToUpper(ch)
	}

	return eventKey == key
}

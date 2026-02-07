package keys

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// ParseKey converts a key string to tcell key, modifiers, and rune.
// Supports formats like:
//   - "Escape", "Enter", "Tab" (special keys)
//   - "Ctrl+Q", "Alt+F1", "Shift+Tab", "Ctrl+Alt+Q" (modified keys)
//   - "a", "Q", "space" (character keys)
func ParseKey(keyStr string) (tcell.Key, tcell.ModMask, rune, error) {
	keyStr = strings.TrimSpace(keyStr)
	if keyStr == "" {
		return 0, 0, 0, fmt.Errorf("empty key string")
	}

	parts := strings.Split(keyStr, "+")

	// Process all but the last part as modifiers
	var modMask tcell.ModMask
	for i := 0; i < len(parts)-1; i++ {
		mod := strings.ToLower(strings.TrimSpace(parts[i]))
		switch mod {
		case "ctrl", "control":
			modMask |= tcell.ModCtrl
		case "alt":
			modMask |= tcell.ModAlt
		case "meta":
			modMask |= tcell.ModMeta
		case "shift":
			modMask |= tcell.ModShift
		default:
			return 0, 0, 0, fmt.Errorf("unknown modifier: %s", parts[i])
		}
	}

	// The last part is the actual key
	keyPart := strings.TrimSpace(parts[len(parts)-1])
	keyLower := strings.ToLower(keyPart)

	// Special keys
	specialKeys := map[string]tcell.Key{
		"escape": tcell.KeyEscape, "esc": tcell.KeyEscape,
		"enter": tcell.KeyEnter, "return": tcell.KeyEnter,
		"tab": tcell.KeyTab, "backtab": tcell.KeyBacktab,
		"backspace": tcell.KeyBackspace, "bs": tcell.KeyBackspace,
		"delete": tcell.KeyDelete, "del": tcell.KeyDelete,
		"insert": tcell.KeyInsert, "ins": tcell.KeyInsert,
		"up": tcell.KeyUp, "down": tcell.KeyDown, "left": tcell.KeyLeft, "right": tcell.KeyRight,
		"home": tcell.KeyHome, "end": tcell.KeyEnd,
		"pgup": tcell.KeyPgUp, "pageup": tcell.KeyPgUp,
		"pgdn": tcell.KeyPgDn, "pagedown": tcell.KeyPgDn,
	}
	if k, ok := specialKeys[keyLower]; ok {
		return k, modMask, 0, nil
	}

	// space key - returns rune ' '
	if keyLower == "space" {
		return tcell.KeyRune, modMask, ' ', nil
	}

	// Function keys F1-F12 (compact parsing)
	if strings.HasPrefix(keyLower, "f") && len(keyPart) > 1 {
		var fNum int
		_, err := fmt.Sscanf(keyPart, "F%d", &fNum)
		if err != nil {
			_, err = fmt.Sscanf(keyPart, "f%d", &fNum)
		}
		if err == nil && fNum >= 1 && fNum <= 12 {
			return tcell.KeyF1 + tcell.Key(fNum-1), modMask, 0, nil
		}
	}

	// Single rune (handles Unicode correctly)
	runes := []rune(keyPart)
	if len(runes) == 1 {
		return tcell.KeyRune, modMask, runes[0], nil
	}

	return 0, 0, 0, fmt.Errorf("unknown key: %s", keyPart)
}

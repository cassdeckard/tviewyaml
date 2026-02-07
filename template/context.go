package template

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Context provides the execution context for templates
type Context struct {
	App    *tview.Application
	Pages  *tview.Pages
	Colors *ColorHelper
}

// NewContext creates a new template context
func NewContext(app *tview.Application, pages *tview.Pages) *Context {
	return &Context{
		App:    app,
		Pages:  pages,
		Colors: &ColorHelper{},
	}
}

// ColorHelper provides color parsing utilities
type ColorHelper struct{}

// Parse converts color names to tcell.Color
func (c *ColorHelper) Parse(name string) tcell.Color {
	colorMap := map[string]tcell.Color{
		"black":  tcell.ColorBlack,
		"red":    tcell.ColorRed,
		"green":  tcell.ColorGreen,
		"yellow": tcell.ColorYellow,
		"blue":   tcell.ColorBlue,
		"white":  tcell.ColorWhite,
		"gray":   tcell.ColorGray,
		"orange": tcell.ColorOrange,
		"purple": tcell.ColorPurple,
		"pink":   tcell.ColorPink,
		"lime":   tcell.ColorLime,
		"navy":   tcell.ColorNavy,
		"teal":   tcell.ColorTeal,
		"silver": tcell.ColorSilver,
	}

	if color, ok := colorMap[name]; ok {
		return color
	}
	return tcell.ColorWhite
}

// ParseAlignment converts alignment strings to tview alignment constants
func ParseAlignment(align string) int {
	switch align {
	case "left":
		return tview.AlignLeft
	case "center":
		return tview.AlignCenter
	case "right":
		return tview.AlignRight
	default:
		return tview.AlignLeft
	}
}

package template

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// BoundView refreshes a view when its state key changes; used for deferred refresh from SetInputCapture.
type BoundView struct {
	Refresh func() string // returns evaluated template string
	SetText func(string)  // applies the string to the view
}

// Context provides the execution context for templates
type Context struct {
	App    *tview.Application
	Pages  *tview.Pages
	Colors *ColorHelper

	state               map[string]interface{}
	subscribers         map[string][]func(interface{})
	boundViews          map[string][]BoundView // key -> views to refresh when key changes
	dirtyKeys           map[string]bool
	formSubmitCallbacks map[string]func() // form name -> callback (e.g. onSubmit)
	formCancelCallbacks map[string]func() // form name -> callback (e.g. onCancel)
	executor            *Executor         // set by app builder so RunCallback can execute templates
	mu                  sync.RWMutex
}

// NewContext creates a new template context
func NewContext(app *tview.Application, pages *tview.Pages) *Context {
	return &Context{
		App:                 app,
		Pages:               pages,
		Colors:              &ColorHelper{},
		state:               make(map[string]interface{}),
		subscribers:         make(map[string][]func(interface{})),
		boundViews:          make(map[string][]BoundView),
		dirtyKeys:           make(map[string]bool),
		formSubmitCallbacks: make(map[string]func()),
		formCancelCallbacks: make(map[string]func()),
	}
}

// SetState updates the view model state and notifies subscribers.
// Safe to call from goroutines; notifications run on main via App.QueueUpdateDraw.
func (c *Context) SetState(key string, value interface{}) {
	if c.App == nil {
		return
	}
	c.App.QueueUpdateDraw(func() {
		c.setStateInternal(key, value)
	})
}

// SetStateDirect updates state and marks the key dirty. Bound views are refreshed
// later from RefreshDirtyBoundViews (e.g. in SetInputCapture) to avoid deadlock.
func (c *Context) SetStateDirect(key string, value interface{}) {
	c.mu.Lock()
	c.state[key] = value
	c.dirtyKeys[key] = true
	c.mu.Unlock()
}

func (c *Context) setStateInternal(key string, value interface{}) {
	c.mu.Lock()
	c.state[key] = value
	c.dirtyKeys[key] = true
	c.mu.Unlock()
}

// RegisterBoundView registers a view that displays state for key. It will be
// refreshed when RefreshDirtyBoundViews runs (on next key event), not from inside handlers.
func (c *Context) RegisterBoundView(key string, bv BoundView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.boundViews[key] = append(c.boundViews[key], bv)
}

// HasDirtyKeys returns true if any state key has been marked dirty (e.g. by SetStateDirect).
func (c *Context) HasDirtyKeys() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.dirtyKeys) > 0
}

// RefreshDirtyBoundViews evaluates and updates all bound views for dirty keys,
// then runs OnStateChange callbacks for those keys. Must be run on the main goroutine
// (e.g. via QueueUpdateDraw from a background refresh goroutine).
func (c *Context) RefreshDirtyBoundViews() {
	c.mu.Lock()
	keys := make([]string, 0, len(c.dirtyKeys))
	for k := range c.dirtyKeys {
		keys = append(keys, k)
		delete(c.dirtyKeys, k)
	}
	viewsByKey := make(map[string][]BoundView)
	subsByKey := make(map[string][]func(interface{}))
	for _, k := range keys {
		viewsByKey[k] = append([]BoundView{}, c.boundViews[k]...)
		subsByKey[k] = append([]func(interface{}){}, c.subscribers[k]...)
	}
	stateCopy := make(map[string]interface{})
	for k, v := range c.state {
		stateCopy[k] = v
	}
	c.mu.Unlock()
	for _, k := range keys {
		for _, bv := range viewsByKey[k] {
			if bv.Refresh != nil && bv.SetText != nil {
				s := bv.Refresh()
				bv.SetText(s)
			}
		}
		for _, fn := range subsByKey[k] {
			if v, ok := stateCopy[k]; ok {
				fn(v)
			}
		}
	}
}

// GetState returns the current value for a state key.
func (c *Context) GetState(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.state[key]
	return v, ok
}

// OnStateChange subscribes to state changes for the given key.
// The callback runs on the main goroutine when SetState is called for that key.
func (c *Context) OnStateChange(key string, fn func(interface{})) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.subscribers[key] = append(c.subscribers[key], fn)
}

// RegisterFormSubmit registers a form's submit callback by name so runFormSubmit(formName) can invoke it (e.g. from a button).
func (c *Context) RegisterFormSubmit(name string, callback func()) {
	if name == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.formSubmitCallbacks[name] = callback
}

// RunFormSubmit runs the submit callback registered for the given form name. No-op if name is unknown.
func (c *Context) RunFormSubmit(name string) {
	c.mu.RLock()
	cb := c.formSubmitCallbacks[name]
	c.mu.RUnlock()
	if cb != nil {
		cb()
	}
}

// RegisterFormCancel registers a form's cancel callback by name so runFormCancel(formName) can invoke it (e.g. from a button).
func (c *Context) RegisterFormCancel(name string, callback func()) {
	if name == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.formCancelCallbacks[name] = callback
}

// RunFormCancel runs the cancel callback registered for the given form name. No-op if name is unknown.
func (c *Context) RunFormCancel(name string) {
	c.mu.RLock()
	cb := c.formCancelCallbacks[name]
	c.mu.RUnlock()
	if cb != nil {
		cb()
	}
}

// SetExecutor sets the template executor so RunCallback can execute template expressions (e.g. from modal onDone).
// Called by the app builder after creating the executor.
func (c *Context) SetExecutor(e *Executor) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.executor = e
}

// RunCallback executes a template expression (e.g. "switchToPage \"main\"") and runs the resulting callback.
// No-op if executor is not set or execution fails.
func (c *Context) RunCallback(templateStr string) {
	c.mu.RLock()
	e := c.executor
	c.mu.RUnlock()
	if e == nil {
		return
	}
	cb, err := e.ExecuteCallback(templateStr)
	if err == nil {
		cb()
	}
}

// ColorHelper provides color parsing utilities
type ColorHelper struct{}

// Parse converts color names to tcell.Color
func (c *ColorHelper) Parse(name string) tcell.Color {
	colorMap := map[string]tcell.Color{
		"black":   tcell.ColorBlack,
		"red":     tcell.ColorRed,
		"green":   tcell.ColorGreen,
		"yellow":  tcell.ColorYellow,
		"blue":    tcell.ColorBlue,
		"white":   tcell.ColorWhite,
		"gray":    tcell.ColorGray,
		"orange":  tcell.ColorOrange,
		"purple":  tcell.ColorPurple,
		"pink":    tcell.ColorPink,
		"lime":    tcell.ColorLime,
		"navy":    tcell.ColorNavy,
		"teal":    tcell.ColorTeal,
		"silver":  tcell.ColorSilver,
		"cyan":    tcell.ColorAqua,    // Cyan/Aqua
		"magenta": tcell.ColorFuchsia, // Magenta/Fuchsia
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

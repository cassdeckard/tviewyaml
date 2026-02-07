package builder

import (
	"github.com/rivo/tview"
)

// CallbackAttacher handles attaching callbacks to primitives
type CallbackAttacher struct{}

// NewCallbackAttacher creates a new callback attacher
func NewCallbackAttacher() *CallbackAttacher {
	return &CallbackAttacher{}
}

// AttachCallback attaches a callback function to a primitive
func (ca *CallbackAttacher) AttachCallback(primitive tview.Primitive, callback func()) error {
	switch v := primitive.(type) {
	case *tview.Button:
		v.SetSelectedFunc(callback)
	case *tview.Checkbox:
		v.SetChangedFunc(func(checked bool) {
			callback()
		})
	// Note: List item callbacks are handled differently during item creation
	default:
		// Some primitives don't have a standard callback mechanism
	}

	return nil
}

// AttachChangeCallback attaches a change callback to a primitive
func (ca *CallbackAttacher) AttachChangeCallback(primitive tview.Primitive, callback func(text string)) error {
	switch v := primitive.(type) {
	case *tview.InputField:
		v.SetChangedFunc(callback)
	default:
		// Not all primitives support change callbacks
	}

	return nil
}

package app

import (
	"sync"
	"time"

	"github.com/cassdeckard/tviewyaml"
	"github.com/cassdeckard/tviewyaml/template"
)

var (
	clockMu   sync.Mutex
	clockDone chan struct{}
)

// RegisterClock adds startClock and stopClock as custom template functions
// so the clock demo works via {{ startClock }} / {{ stopClock }} in YAML.
// All clock logic lives in the example, not in the core library.
func RegisterClock(b *tviewyaml.AppBuilder) *tviewyaml.AppBuilder {
	maxZero := 0
	b.WithTemplateFunction("startClock", 0, &maxZero, nil, startClock)
	b.WithTemplateFunction("stopClock", 0, &maxZero, nil, stopClock)
	return b
}

func startClock(ctx *template.Context) {
	clockMu.Lock()
	if clockDone != nil {
		close(clockDone)
		clockDone = nil
	}
	clockDone = make(chan struct{})
	done := clockDone
	clockMu.Unlock()

	// Show time immediately (we're on main goroutine in button handler)
	ctx.SetStateDirect("clock", time.Now().Format("15:04:05"))

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				app := ctx.App
				if app == nil {
					return
				}
				app.QueueUpdateDraw(func() {
					ctx.SetStateDirect("clock", t.Format("15:04:05"))
				})
			}
		}
	}()
}

func stopClock(ctx *template.Context) {
	clockMu.Lock()
	defer clockMu.Unlock()
	if clockDone != nil {
		close(clockDone)
		clockDone = nil
	}
}

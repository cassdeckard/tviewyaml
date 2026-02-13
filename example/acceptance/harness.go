package acceptance_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"example/app"
	"github.com/cassdeckard/tviewyaml"
	"github.com/cassdeckard/tviewyaml/keys"
	"github.com/gdamore/tcell/v2"
)

const drawTimeout = 3 * time.Second

// acceptanceHarness runs the example app with a SimulationScreen and provides
// helpers to wait for draws, inject input, and assert on content.
type acceptanceHarness struct {
	app       *tviewyaml.Application
	drawDone  chan struct{}
	contentMu sync.Mutex
	content   string
	runDone   chan struct{}
}

// newAcceptanceHarness builds the example app with a simulation screen at the given size,
// starts Run() in a goroutine, and sets up draw synchronization. Caller must call stop() when done.
func newAcceptanceHarness(t *testing.T, cols, rows int) *acceptanceHarness {
	t.Helper()
	sim := tcell.NewSimulationScreen("UTF-8")
	if err := sim.Init(); err != nil {
		t.Fatalf("SimulationScreen Init: %v", err)
	}
	// SimulationScreen does not expose SetSize in the public interface; we inject EventResize after start.

	application, pageErrors, err := app.BuildWithScreen("../config", sim)
	if err != nil {
		sim.Fini()
		t.Fatalf("Build: %v", err)
	}
	if len(pageErrors) > 0 {
		t.Logf("Build had %d page errors (non-fatal): %v", len(pageErrors), pageErrors)
	}

	drawDone := make(chan struct{}, 1)
	h := &acceptanceHarness{app: application, drawDone: drawDone, runDone: make(chan struct{})}

	application.SetAfterDrawFunc(func(screen tcell.Screen) {
		w, hi := screen.Size()
		var b strings.Builder
		for y := 0; y < hi; y++ {
			for x := 0; x < w; x++ {
				mainc, _, _, _ := screen.GetContent(x, y)
				if mainc != 0 {
					b.WriteRune(mainc)
				} else {
					b.WriteByte(' ')
				}
			}
			if y < hi-1 {
				b.WriteByte('\n')
			}
		}
		h.contentMu.Lock()
		h.content = b.String()
		h.contentMu.Unlock()
		select {
		case h.drawDone <- struct{}{}:
		default:
			// already one pending
		}
	})

	go func() {
		defer close(h.runDone)
		_ = application.Run()
	}()

	// Trigger initial resize so layout runs at desired size
	h.resize(cols, rows)
	if !h.waitForDraw() {
		application.Stop()
		<-h.runDone
		sim.Fini()
		t.Fatal("timeout waiting for initial draw")
	}
	return h
}

func (h *acceptanceHarness) waitForDraw() bool {
	select {
	case <-h.drawDone:
		return true
	case <-time.After(drawTimeout):
		return false
	}
}

// waitForDraws waits for n draws (use after injecting input to see the resulting screen).
func (h *acceptanceHarness) waitForDraws(n int) bool {
	for i := 0; i < n; i++ {
		if !h.waitForDraw() {
			return false
		}
	}
	return true
}

func (h *acceptanceHarness) getContent() string {
	h.contentMu.Lock()
	defer h.contentMu.Unlock()
	return h.content
}

func (h *acceptanceHarness) screenContains(substr string) bool {
	return strings.Contains(h.getContent(), substr)
}

func (h *acceptanceHarness) resize(cols, rows int) {
	h.app.QueueEvent(tcell.NewEventResize(cols, rows))
}

func (h *acceptanceHarness) typeKey(keyStr string) {
	tcellKey, mod, r, err := keys.ParseKey(keyStr)
	if err != nil {
		panic("typeKey: " + err.Error())
	}
	h.app.QueueEvent(tcell.NewEventKey(tcellKey, r, mod))
}

func (h *acceptanceHarness) stop() {
	h.app.Stop()
	<-h.runDone
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

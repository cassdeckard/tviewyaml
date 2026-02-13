package acceptance_test

import (
	"fmt"
	"os"
	"path/filepath"
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

const snapshotEnvUpdate = "UPDATE_TERMINAL_SNAPSHOTS"

// terminalSizes are common sizes used for multi-size snapshot tests.
var terminalSizes = []struct {
	name       string
	cols, rows int
}{
	{"80x24", 80, 24},
	{"120x30", 120, 30},
	{"40x10", 40, 10},
}

// runAtSizes runs fn as a subtest for each terminal size. Each subtest gets its own harness.
func runAtSizes(t *testing.T, fn func(t *testing.T, h *acceptanceHarness)) {
	t.Helper()
	for _, sz := range terminalSizes {
		sz := sz
		t.Run(sz.name, func(t *testing.T) {
			t.Helper()
			h := newAcceptanceHarness(t, sz.cols, sz.rows)
			defer h.stop()
			fn(t, h)
		})
	}
}

// TerminalSnapshot is a point-in-time capture of the simulated terminal (character grid and dimensions).
// Content is newline-separated lines; String() returns Content so it can be echoed or logged.
type TerminalSnapshot struct {
	Content string
	Cols   int
	Rows   int
}

// String returns the terminal content so that t.Log(snap) or echo displays the terminal.
// In a narrower real terminal, long lines wrap naturally.
func (s TerminalSnapshot) String() string {
	return s.Content
}

// DelimitedString returns the snapshot with a header and footer for extraction from test output.
func (s TerminalSnapshot) DelimitedString() string {
	return "--- terminal snapshot " + fmt.Sprintf("%dx%d", s.Cols, s.Rows) + " ---\n" +
		s.Content + "\n--- end snapshot ---"
}

// acceptanceHarness runs the example app with a SimulationScreen and provides
// helpers to wait for draws, inject input, and assert on content.
type acceptanceHarness struct {
	app       *tviewyaml.Application
	drawDone  chan struct{}
	contentMu sync.Mutex
	content   string
	lastCols  int
	lastRows  int
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
	application, pageErrors, err := app.BuildWithScreen("../config", sim)
	if err != nil {
		sim.Fini()
		t.Fatalf("Build: %v", err)
	}
	if len(pageErrors) > 0 {
		t.Logf("Build had %d page errors (non-fatal): %v", len(pageErrors), pageErrors)
	}

	// Set simulation screen size so the first draw uses the correct dimensions.
	sim.SetSize(cols, rows)

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
		h.lastCols = w
		h.lastRows = hi
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

	// Queue resize so the app's layout runs at desired size; screen is already SetSize'd above.
	h.resize(cols, rows)
	if !h.waitForDraw() {
		application.Stop()
		<-h.runDone
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

// TakeSnapshot returns the current terminal content and dimensions.
// Call waitForDraw() first if a fresh frame is needed.
func (h *acceptanceHarness) TakeSnapshot() TerminalSnapshot {
	h.contentMu.Lock()
	defer h.contentMu.Unlock()
	return TerminalSnapshot{Content: h.content, Cols: h.lastCols, Rows: h.lastRows}
}

// snapshotGoldenPath returns the path to the golden file for the given name.
// Name is sanitized for use as a filename (e.g. subtest "80x24" -> "TestName_80x24.terminal" when name is derived from t.Name()).
func snapshotGoldenPath(name string) string {
	safe := strings.ReplaceAll(name, "/", "_")
	safe = strings.TrimSpace(safe)
	if safe == "" {
		safe = "default"
	}
	return filepath.Join("testdata", "snapshots", safe+".terminal")
}

// AssertSnapshot compares the current terminal state to the golden snapshot at testdata/snapshots/<name>.terminal.
// If name is empty, the name is derived from t.Name() (e.g. TestAcceptance_Layout/80x24 -> TestAcceptance_Layout_80x24.terminal).
// When UPDATE_TERMINAL_SNAPSHOTS=1 is set, the golden file is overwritten with the current state and the assertion passes.
func (h *acceptanceHarness) AssertSnapshot(t *testing.T, name string) {
	t.Helper()
	if name == "" {
		name = strings.ReplaceAll(t.Name(), "/", "_")
	}
	snap := h.TakeSnapshot()
	path := snapshotGoldenPath(name)

	if os.Getenv(snapshotEnvUpdate) != "" {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("create snapshot dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(snap.Content), 0644); err != nil {
			t.Fatalf("write snapshot: %v", err)
		}
		return
	}

	expected, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Logf("current terminal:\n%s", snap)
			t.Fatalf("no golden snapshot at %s; run with UPDATE_TERMINAL_SNAPSHOTS=1 to create it", path)
		}
		t.Fatalf("read golden snapshot: %v", err)
	}
	expectedStr := string(expected)
	if snap.Content != expectedStr {
		t.Errorf("snapshot mismatch for %s", name)
		t.Logf("current terminal:\n%s", snap)
		t.Logf("expected (golden):\n%s", expectedStr)
	}
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

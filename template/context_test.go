package template

import (
	"sync"
	"testing"
)

// TestContextConcurrency verifies that concurrent reads and writes to context state
// are safe and don't cause data races. This test should be run with the race detector:
// go test -race ./template
func TestContextConcurrency(t *testing.T) {
	ctx := newTestContext()
	done := make(chan bool)
	iterations := 1000

	// Writer goroutine - continuously updates state
	go func() {
		for i := 0; i < iterations; i++ {
			ctx.SetStateDirect("key", i)
		}
		done <- true
	}()

	// Reader goroutine - continuously reads state
	go func() {
		for i := 0; i < iterations; i++ {
			_, _ = ctx.GetState("key")
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done
}

// TestContextConcurrentMultipleKeys tests concurrent access to multiple different keys
func TestContextConcurrentMultipleKeys(t *testing.T) {
	ctx := newTestContext()
	var wg sync.WaitGroup
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	iterations := 200

	// Start a writer and reader for each key
	for _, key := range keys {
		key := key // capture for goroutine
		wg.Add(2)

		// Writer for this key
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				ctx.SetStateDirect(key, i)
			}
		}()

		// Reader for this key
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_, _ = ctx.GetState(key)
			}
		}()
	}

	wg.Wait()
}

// TestContextConcurrentDirtyTracking tests concurrent access to dirty key tracking
func TestContextConcurrentDirtyTracking(t *testing.T) {
	ctx := newTestContext()
	var wg sync.WaitGroup
	iterations := 500

	// Multiple goroutines setting state (which marks keys dirty)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				ctx.SetStateDirect("key", id*iterations+j)
			}
		}(i)
	}

	// Goroutine checking dirty status
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = ctx.HasDirtyKeys()
		}
	}()

	wg.Wait()
}

// TestContextConcurrentBoundViews tests concurrent access to bound view registration
func TestContextConcurrentBoundViews(t *testing.T) {
	ctx := newTestContext()
	var wg sync.WaitGroup
	iterations := 100

	// Create a bound view
	bv := BoundView{
		Refresh: func() string { return "test" },
		SetText: func(s string) {},
	}

	// Multiple goroutines registering bound views
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				ctx.RegisterBoundView("key", bv)
			}
		}(i)
	}

	wg.Wait()
}

// TestContextConcurrentSubscribers tests concurrent access to state change subscribers
func TestContextConcurrentSubscribers(t *testing.T) {
	ctx := newTestContext()
	var wg sync.WaitGroup
	iterations := 100

	// Multiple goroutines registering subscribers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				ctx.OnStateChange("key", func(v interface{}) {})
			}
		}(i)
	}

	wg.Wait()
}

// TestContextConcurrentFormCallbacks tests concurrent access to form callbacks
func TestContextConcurrentFormCallbacks(t *testing.T) {
	ctx := newTestContext()
	var wg sync.WaitGroup
	iterations := 200

	// Writer goroutines registering callbacks
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ctx.RegisterFormSubmit("form1", func() {})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ctx.RegisterFormSubmit("form2", func() {})
		}
	}()

	// Reader goroutines calling callbacks
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ctx.RunFormSubmit("form1")
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ctx.RunFormSubmit("form2")
		}
	}()

	wg.Wait()
}

// TestContextConcurrentExecutor tests concurrent access to executor
func TestContextConcurrentExecutor(t *testing.T) {
	ctx := newTestContext()
	executor, _ := newTestExecutor()
	var wg sync.WaitGroup
	iterations := 100

	// Writer setting executor
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			ctx.SetExecutor(executor)
		}
	}()

	// Readers calling RunCallback
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				ctx.RunCallback("{{ testFunc }}")
			}
		}()
	}

	wg.Wait()
}

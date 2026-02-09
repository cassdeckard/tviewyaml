package main

import (
	"fmt"
	"sync/atomic"

	"github.com/cassdeckard/tviewyaml"
	"github.com/cassdeckard/tviewyaml/template"
)

var (
	counter       int32
	messageIndex  int32
	messages      = []string{
		"Hello, World!",
		"State binding is powerful!",
		"All views update together!",
		"Reactive UI made easy!",
		"Try clicking again!",
	}
)

// RegisterStateBinding adds state binding demo functions as custom template functions
// so the state-binding demo works via {{ incrementCounter }}, {{ resetCounter }}, etc.
func RegisterStateBinding(b *tviewyaml.AppBuilder) *tviewyaml.AppBuilder {
	maxZero := 0
	b.WithTemplateFunction("incrementCounter", 0, &maxZero, nil, incrementCounter)
	b.WithTemplateFunction("resetCounter", 0, &maxZero, nil, resetCounter)
	b.WithTemplateFunction("updateMessage", 0, &maxZero, nil, updateMessage)
	return b
}

func incrementCounter(ctx *template.Context) {
	// Increment counter atomically
	newValue := atomic.AddInt32(&counter, 1)
	
	// Update state - this will trigger all bound views to refresh
	ctx.SetStateDirect("counter", fmt.Sprintf("Count: %d", newValue))
}

func resetCounter(ctx *template.Context) {
	// Reset counter to 0
	atomic.StoreInt32(&counter, 0)
	
	// Update state - this will trigger all bound views to refresh
	ctx.SetStateDirect("counter", "Count: 0")
}

func updateMessage(ctx *template.Context) {
	// Get next message index (cycle through messages)
	idx := atomic.AddInt32(&messageIndex, 1) % int32(len(messages))
	
	// Update state with new message
	ctx.SetStateDirect("message", messages[idx])
}

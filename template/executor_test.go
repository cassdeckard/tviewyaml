package template

import (
	"strings"
	"testing"

	"github.com/rivo/tview"
)

// newTestContext creates a Context for testing
func newTestContext() *Context {
	app := tview.NewApplication()
	pages := tview.NewPages()
	return NewContext(app, pages)
}

// newTestRegistry creates a registry with test evaluators and functions
func newTestRegistry() *FunctionRegistry {
	registry := NewFunctionRegistry()

	// Register test evaluators
	registry.RegisterEvaluator("testEval", 1, 1, func(ctx *Context, args []string) string {
		return "eval:" + args[0]
	})

	registry.RegisterEvaluator("testEvalNoArgs", 0, 0, func(ctx *Context, args []string) string {
		return "noargs"
	})

	registry.RegisterEvaluator("bindState", 1, 1, func(ctx *Context, args []string) string {
		v, ok := ctx.GetState(args[0])
		if !ok {
			return ""
		}
		return "state:" + args[0] + "=" + strings.TrimSpace(v.(string))
	})

	// Register test functions
	intPtr := func(i int) *int { return &i }

	registry.Register("testFunc", 0, intPtr(0), nil, func(ctx *Context) {
		ctx.SetStateDirect("testFunc", "called")
	})

	registry.Register("testFuncOneArg", 1, intPtr(1), nil, func(ctx *Context, arg1 string) {
		ctx.SetStateDirect("testFuncOneArg", arg1)
	})

	registry.Register("testFuncTwoArgs", 2, intPtr(2), nil, func(ctx *Context, arg1, arg2 string) {
		ctx.SetStateDirect("testFuncTwoArgs", arg1+"|"+arg2)
	})

	registry.Register("testFuncVariadic", 1, nil, nil, func(ctx *Context, args []string) {
		ctx.SetStateDirect("testFuncVariadic", strings.Join(args, ","))
	})

	registry.Register("testFuncWithValidator", 1, intPtr(1), func(ctx *Context, args []string) error {
		if args[0] == "invalid" {
			return &testValidationError{msg: "validation failed"}
		}
		return nil
	}, func(ctx *Context, arg1 string) {
		ctx.SetStateDirect("testFuncWithValidator", arg1)
	})

	return registry
}

type testValidationError struct {
	msg string
}

func (e *testValidationError) Error() string {
	return e.msg
}

// newTestExecutor creates an Executor for testing
func newTestExecutor() (*Executor, *Context) {
	ctx := newTestContext()
	registry := newTestRegistry()
	executor := NewExecutor(ctx, registry)
	return executor, ctx
}

func TestEvaluateToString(t *testing.T) {
	executor, ctx := newTestExecutor()

	tests := []struct {
		name        string
		templateStr string
		setupState  func() // optional setup
		want        string
		wantErr     bool
		errContains string
	}{
		// Empty and plain text
		{"empty string", "", nil, "", false, ""},
		{"plain text", "Hello World", nil, "Hello World", false, ""},
		{"plain text with spaces", "  Hello   World  ", nil, "  Hello   World  ", false, ""},

		// Single evaluator
		{"single evaluator", "{{ testEval hello }}", nil, "eval:hello", false, ""},
		{"single evaluator with spaces", "{{ testEvalNoArgs }}", nil, "noargs", false, ""},

		// Multiple evaluators
		{"two evaluators", "{{ testEval a }} {{ testEval b }}", nil, "eval:a eval:b", false, ""},
		{"multiple evaluators", "{{ testEval 1 }} {{ testEval 2 }} {{ testEval 3 }}", nil, "eval:1 eval:2 eval:3", false, ""},

		// Mixed literal and evaluators
		{"prefix literal", "Hello {{ testEval world }}", nil, "Hello eval:world", false, ""},
		{"suffix literal", "{{ testEval hello }} World", nil, "eval:hello World", false, ""},
		{"both sides", "Hello {{ testEval world }}!", nil, "Hello eval:world!", false, ""},
		{"multiple mixed", "A {{ testEval 1 }} B {{ testEval 2 }} C", nil, "A eval:1 B eval:2 C", false, ""},

		// bindState evaluator
		{"bindState exists", "{{ bindState key1 }}", func() {
			ctx.SetStateDirect("key1", "value1")
		}, "value1", false, ""},
		{"bindState missing", "{{ bindState missing }}", nil, "", false, ""},
		{"bindState multiple", "{{ bindState a }} {{ bindState b }}", func() {
			ctx.SetStateDirect("a", "A")
			ctx.SetStateDirect("b", "B")
		}, "A B", false, ""},

		// Error cases
		{"unknown evaluator", "{{ unknownEval }}", nil, "", true, "unknown evaluator"},
		{"wrong arg count too few", "{{ testEval }}", nil, "", true, "expects 1-1 args"},
		{"wrong arg count too many", "{{ testEvalNoArgs extra }}", nil, "", true, "expects 0-0 args"},

		// Edge cases
		{"unclosed template", "{{ testEval hello", nil, "", true, "unknown evaluator"}, // unclosed template causes parsing issues
		{"empty template block", "{{ }}", nil, "", true, "unknown evaluator"},
		{"whitespace only template", "{{   }}", nil, "", true, "unknown evaluator"},
		{"nested braces literal", "{{ testEval {hello} }}", nil, "eval:{hello}", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset state
			ctx.mu.Lock()
			ctx.state = make(map[string]interface{})
			ctx.mu.Unlock()

			if tt.setupState != nil {
				tt.setupState()
			}

			got, err := executor.EvaluateToString(tt.templateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateToString(%q) error = %v, wantErr %v", tt.templateStr, err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("EvaluateToString(%q) error = %v, want error containing %q", tt.templateStr, err, tt.errContains)
					}
				}
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateToString(%q) = %q, want %q", tt.templateStr, got, tt.want)
			}
		})
	}
}

func TestExecuteCallback(t *testing.T) {
	executor, ctx := newTestExecutor()

	tests := []struct {
		name        string
		templateStr string
		wantErr     bool
		errContains string
		verify      func() bool // verify callback execution
	}{
		// Empty callback
		{"empty string", "", false, "", func() bool { return true }}, // no-op callback

		// Function calls with no args
		{"no args", "{{ testFunc }}", false, "", func() bool {
			v, ok := ctx.GetState("testFunc")
			return ok && v == "called"
		}},

		// Function calls with args
		{"one arg", "{{ testFuncOneArg \"hello\" }}", false, "", func() bool {
			v, ok := ctx.GetState("testFuncOneArg")
			return ok && v == "hello"
		}},
		{"two args", "{{ testFuncTwoArgs \"a\" \"b\" }}", false, "", func() bool {
			v, ok := ctx.GetState("testFuncTwoArgs")
			return ok && v == "a|b"
		}},
		{"three args variadic", "{{ testFuncVariadic \"a\" \"b\" \"c\" }}", false, "", func() bool {
			v, ok := ctx.GetState("testFuncVariadic")
			return ok && v == "a,b,c"
		}},

		// Quoted arguments
		{"quoted with spaces", "{{ testFuncOneArg \"hello world\" }}", false, "", func() bool {
			v, ok := ctx.GetState("testFuncOneArg")
			return ok && v == "hello world"
		}},
		{"quoted with escaped quotes", "{{ testFuncOneArg \"hello\\\"world\" }}", false, "", func() bool {
			v, ok := ctx.GetState("testFuncOneArg")
			return ok && v == "hello\"world"
		}},

		// Error cases
		{"unknown function", "{{ unknownFunc }}", true, "unknown function", nil},
		{"too few args", "{{ testFuncOneArg }}", true, "requires at least 1 argument", nil},
		{"too many args", "{{ testFunc \"extra\" }}", true, "accepts at most 0 argument", nil},
		{"validator failure", "{{ testFuncWithValidator \"invalid\" }}", true, "validation failed", nil},
		{"validator success", "{{ testFuncWithValidator \"valid\" }}", false, "", func() bool {
			v, ok := ctx.GetState("testFuncWithValidator")
			return ok && v == "valid"
		}},

		// Edge cases
		{"invalid expression", "not a template", true, "unknown function", nil}, // "not" is parsed as function name
		{"missing closing brace", "{{ testFunc", false, "", func() bool {
			// Missing closing brace - ExecuteCallback strips {{ and }}, so this becomes "testFunc"
			// which is a valid function call
			v, ok := ctx.GetState("testFunc")
			return ok && v == "called"
		}},
		{"empty function name", "{{ \"arg\" }}", true, "invalid template expression", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset state
			ctx.mu.Lock()
			ctx.state = make(map[string]interface{})
			ctx.mu.Unlock()

			callback, err := executor.ExecuteCallback(tt.templateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteCallback(%q) error = %v, wantErr %v", tt.templateStr, err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("ExecuteCallback(%q) error = %v, want error containing %q", tt.templateStr, err, tt.errContains)
					}
				}
				return
			}

			// Execute callback and verify
			if callback != nil {
				callback()
			}

			if tt.verify != nil && !tt.verify() {
				t.Errorf("ExecuteCallback(%q) callback execution verification failed", tt.templateStr)
			}
		})
	}
}

func TestExtractBindStateKeys(t *testing.T) {
	executor, _ := newTestExecutor()

	tests := []struct {
		name        string
		templateStr string
		want        []string
	}{
		{"empty string", "", nil},
		{"no bindState", "Hello World", nil},
		{"no bindState with other evaluator", "{{ testEval hello }}", nil},
		{"single bindState", "{{ bindState key1 }}", []string{"key1"}},
		{"single bindState with literal", "Hello {{ bindState key1 }}", []string{"key1"}},
		{"multiple bindState different keys", "{{ bindState a }} {{ bindState b }}", []string{"a", "b"}},
		{"multiple bindState same key", "{{ bindState key1 }} {{ bindState key1 }}", []string{"key1"}}, // deduplicated
		{"multiple bindState mixed", "{{ bindState a }} {{ testEval x }} {{ bindState b }} {{ bindState a }}", []string{"a", "b"}},
		{"bindState with spaces", "{{ bindState  key1  }}", []string{"key1"}},
		{"bindState no args", "{{ bindState }}", nil}, // no args, so not extracted
		{"bindState multiple in complex template", "A {{ bindState x }} B {{ bindState y }} C {{ bindState z }} D", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executor.ExtractBindStateKeys(tt.templateStr)

			// Compare slices (order may vary, so check length and membership)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractBindStateKeys(%q) = %v, want %v", tt.templateStr, got, tt.want)
				return
			}

			// Create maps for comparison
			gotMap := make(map[string]bool)
			for _, k := range got {
				gotMap[k] = true
			}
			wantMap := make(map[string]bool)
			for _, k := range tt.want {
				wantMap[k] = true
			}

			for k := range gotMap {
				if !wantMap[k] {
					t.Errorf("ExtractBindStateKeys(%q) returned unexpected key %q", tt.templateStr, k)
				}
			}
			for k := range wantMap {
				if !gotMap[k] {
					t.Errorf("ExtractBindStateKeys(%q) missing expected key %q", tt.templateStr, k)
				}
			}
		})
	}
}

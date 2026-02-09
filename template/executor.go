package template

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Executor handles template execution
type Executor struct {
	ctx      *Context
	registry *FunctionRegistry
}

// NewExecutor creates a new template executor
func NewExecutor(ctx *Context, registry *FunctionRegistry) *Executor {
	return &Executor{
		ctx:      ctx,
		registry: registry,
	}
}

// EvaluateToString evaluates a template string containing {{ bindState key }} (and other evaluators) and returns the rendered string.
// Example: "Notification: {{ bindState notification }}" -> "Notification: Hello" when state "notification" is "Hello"
func (e *Executor) EvaluateToString(templateStr string) (string, error) {
	if templateStr == "" {
		return "", nil
	}
	return e.evaluateTemplateString(templateStr)
}

// ExtractBindStateKeys returns all state keys referenced by bindState in the template string.
// Used to subscribe to state changes for re-evaluation.
func (e *Executor) ExtractBindStateKeys(templateStr string) []string {
	var keys []string
	seen := make(map[string]bool)
	for _, expr := range extractTemplateExpressions(templateStr) {
		name, args := parseEvaluatorExpr(expr)
		if name == "bindState" && len(args) > 0 && !seen[args[0]] {
			keys = append(keys, args[0])
			seen[args[0]] = true
		}
	}
	return keys
}

// evaluateTemplateString parses {{ ... }} blocks and evaluates them
func (e *Executor) evaluateTemplateString(s string) (string, error) {
	parts := splitTemplateString(s)
	// Pre-allocate buffer capacity: original string length + estimated expansion for evaluators
	estimatedSize := len(s) + len(parts)*16
	var result strings.Builder
	result.Grow(estimatedSize)
	
	for i, part := range parts {
		if i%2 == 0 {
			result.WriteString(part)
			continue
		}
		expr := strings.TrimSpace(part)
		name, args := parseEvaluatorExpr(expr)
		ev, ok := e.registry.GetEvaluator(name)
		if !ok {
			return "", fmt.Errorf("unknown evaluator: %s", name)
		}
		if len(args) < ev.MinArgs || len(args) > ev.MaxArgs {
			return "", fmt.Errorf("evaluator %q expects %d-%d args, got %d", name, ev.MinArgs, ev.MaxArgs, len(args))
		}
		result.WriteString(ev.Handler(e.ctx, args))
	}
	return result.String(), nil
}

// splitTemplateString splits by {{ and }}; even indices are literal, odd are expression content
func splitTemplateString(s string) []string {
	var parts []string
	for {
		start := strings.Index(s, "{{")
		if start < 0 {
			parts = append(parts, s)
			break
		}
		parts = append(parts, s[:start])
		s = s[start+2:]
		end := strings.Index(s, "}}")
		if end < 0 {
			parts = append(parts, "{{"+s) // treat as literal if unclosed
			break
		}
		parts = append(parts, s[:end])
		s = s[end+2:]
	}
	return parts
}

// extractTemplateExpressions returns the content of each {{ ... }} block
func extractTemplateExpressions(s string) []string {
	var exprs []string
	for {
		start := strings.Index(s, "{{")
		if start < 0 {
			break
		}
		s = s[start+2:]
		end := strings.Index(s, "}}")
		if end < 0 {
			break
		}
		exprs = append(exprs, strings.TrimSpace(s[:end]))
		s = s[end+2:]
	}
	return exprs
}

// parseEvaluatorExpr parses "funcName arg1 arg2" into name and args (supports unquoted identifiers)
func parseEvaluatorExpr(expr string) (string, []string) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return "", nil
	}
	re := regexp.MustCompile(`^(\w+)\s*(.*)$`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) < 2 {
		return "", nil
	}
	name := matches[1]
	rest := strings.TrimSpace(matches[2])
	if rest == "" {
		return name, nil
	}
	// Try quoted args first; if none, use unquoted words
	args := parseArguments(rest)
	if len(args) == 0 {
		args = strings.Fields(rest)
	}
	return name, args
}

// ExecuteCallback parses and executes a template expression to create a callback function
func (e *Executor) ExecuteCallback(templateStr string) (func(), error) {
	// Parse the template string to extract function calls
	// Template format: {{ functionName "arg1" "arg2" }}

	if templateStr == "" {
		// Empty template string - return no-op callback
		return func() {}, nil
	}

	// Remove template delimiters
	templateStr = strings.TrimSpace(templateStr)
	templateStr = strings.TrimPrefix(templateStr, "{{")
	templateStr = strings.TrimSuffix(templateStr, "}}")
	templateStr = strings.TrimSpace(templateStr)

	// Parse function name and arguments
	return e.parseAndCreateCallback(templateStr)
}

// parseAndCreateCallback parses the template string and creates the appropriate callback
func (e *Executor) parseAndCreateCallback(expr string) (func(), error) {
	// Match function calls with arguments
	// Pattern: functionName "arg1" "arg2" ...
	re := regexp.MustCompile(`^(\w+)\s*(.*)$`)
	matches := re.FindStringSubmatch(expr)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid template expression: %s", expr)
	}

	funcName := matches[1]
	argsStr := strings.TrimSpace(matches[2])

	// Parse arguments (strings in quotes)
	args := parseArguments(argsStr)

	// Look up function in registry
	fn, ok := e.registry.Get(funcName)
	if !ok {
		return nil, fmt.Errorf("unknown function: %s", funcName)
	}

	// Validate argument count
	if len(args) < fn.MinArgs {
		return nil, fmt.Errorf("function %q requires at least %d argument(s), got %d", funcName, fn.MinArgs, len(args))
	}
	if fn.MaxArgs != nil && len(args) > *fn.MaxArgs {
		return nil, fmt.Errorf("function %q accepts at most %d argument(s), got %d", funcName, *fn.MaxArgs, len(args))
	}

	// Call validator if present (only called after argument count validation)
	if fn.Validator != nil {
		if err := fn.Validator(e.ctx, args); err != nil {
			return nil, fmt.Errorf("validation failed for function %q: %w", funcName, err)
		}
	}

	// Create callback that invokes the handler
	return e.createCallbackFromHandler(fn, args)
}

// parseArguments extracts string arguments from a function call
func parseArguments(argsStr string) []string {
	if argsStr == "" {
		return []string{}
	}

	var args []string
	var current strings.Builder
	inQuote := false
	escaped := false

	for i := 0; i < len(argsStr); i++ {
		ch := argsStr[i]

		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if ch == '"' {
			if inQuote {
				// End of quoted string
				args = append(args, current.String())
				current.Reset()
				inQuote = false
			} else {
				// Start of quoted string
				inQuote = true
			}
			continue
		}

		if inQuote {
			current.WriteByte(ch)
		}
	}

	return args
}

// createCallbackFromHandler creates a callback function that invokes the handler with proper arguments
func (e *Executor) createCallbackFromHandler(fn *TemplateFunction, args []string) (func(), error) {
	handlerValue := reflect.ValueOf(fn.Handler)
	contextValue := reflect.ValueOf(e.ctx)

	return func() {
		// Prepare arguments for the handler call
		var callArgs []reflect.Value
		callArgs = append(callArgs, contextValue)

		if fn.MaxArgs != nil {
			// Fixed args: pass each argument individually
			for _, arg := range args {
				callArgs = append(callArgs, reflect.ValueOf(arg))
			}
		} else {
			// Variadic: pass args as a slice
			callArgs = append(callArgs, reflect.ValueOf(args))
		}

		// Call the handler
		handlerValue.Call(callArgs)
	}, nil
}

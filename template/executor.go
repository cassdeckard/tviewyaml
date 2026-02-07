package template

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rivo/tview"
)

// Executor handles template execution
type Executor struct {
	ctx *Context
}

// NewExecutor creates a new template executor
func NewExecutor(ctx *Context) *Executor {
	return &Executor{
		ctx: ctx,
	}
}

// ExecuteCallback parses and executes a template expression to create a callback function
func (e *Executor) ExecuteCallback(templateStr string) (func(), error) {
	// Parse the template string to extract function calls
	// Template format: {{ functionName "arg1" "arg2" }}

	if templateStr == "" {
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

	// Execute the appropriate function
	switch funcName {
	case "switchToPage":
		if len(args) < 1 {
			return nil, fmt.Errorf("switchToPage requires 1 argument")
		}
		return func() { e.ctx.Pages.SwitchToPage(args[0]) }, nil

	case "removePage":
		if len(args) < 1 {
			return nil, fmt.Errorf("removePage requires 1 argument")
		}
		return func() { e.ctx.Pages.RemovePage(args[0]) }, nil

	case "stopApp":
		return func() { e.ctx.App.Stop() }, nil

	case "showSimpleModal":
		if len(args) < 1 {
			return nil, fmt.Errorf("showSimpleModal requires at least 1 argument (text)")
		}
		text := args[0]
		// Remaining args are button labels
		buttons := args[1:]
		if len(buttons) == 0 {
			buttons = []string{"OK"}
		}
		return func() {
			modal := tview.NewModal().
				SetText(text).
				AddButtons(buttons).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					e.ctx.Pages.RemovePage("simple-modal")
				})
			e.ctx.Pages.AddPage("simple-modal", modal, false, true)
		}, nil

	case "noop":
		return func() {}, nil

	default:
		return nil, fmt.Errorf("unknown function: %s", funcName)
	}
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

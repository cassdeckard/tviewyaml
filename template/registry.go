package template

import (
	"fmt"
	"reflect"
)

// TemplateFunction defines a registered template function
type TemplateFunction struct {
	Name      string
	MinArgs   int
	MaxArgs   *int // nil means unlimited (variadic)
	Validator func(*Context, []string) error
	Handler   interface{} // Function that executes the template logic
}

// TemplateEvaluator defines a value-returning template function (e.g. bindState)
type TemplateEvaluator struct {
	Name    string
	MinArgs int
	MaxArgs int // evaluators use fixed arg count
	Handler func(*Context, []string) string
}

// FunctionRegistry manages registered template functions
type FunctionRegistry struct {
	functions  map[string]*TemplateFunction
	evaluators map[string]*TemplateEvaluator
}

// NewFunctionRegistry creates a new function registry with built-in functions
func NewFunctionRegistry() *FunctionRegistry {
	registry := &FunctionRegistry{
		functions:  make(map[string]*TemplateFunction),
		evaluators: make(map[string]*TemplateEvaluator),
	}
	registerBuiltinFunctions(registry)
	return registry
}

// RegisterEvaluator adds a value-returning template function (e.g. bindState)
func (r *FunctionRegistry) RegisterEvaluator(name string, minArgs, maxArgs int, handler func(*Context, []string) string) error {
	if minArgs < 0 || maxArgs < minArgs {
		return fmt.Errorf("invalid evaluator args: minArgs=%d maxArgs=%d", minArgs, maxArgs)
	}
	if _, exists := r.evaluators[name]; exists {
		return fmt.Errorf("evaluator %q is already registered", name)
	}
	r.evaluators[name] = &TemplateEvaluator{
		Name:    name,
		MinArgs: minArgs,
		MaxArgs: maxArgs,
		Handler: handler,
	}
	return nil
}

// GetEvaluator retrieves an evaluator by name
func (r *FunctionRegistry) GetEvaluator(name string) (*TemplateEvaluator, bool) {
	ev, ok := r.evaluators[name]
	return ev, ok
}

// Register adds a new template function to the registry
func (r *FunctionRegistry) Register(name string, minArgs int, maxArgs *int, validator func(*Context, []string) error, handler interface{}) error {
	// Validate minArgs
	if minArgs < 0 {
		return fmt.Errorf("minArgs must be non-negative, got %d", minArgs)
	}

	// Validate maxArgs if provided
	if maxArgs != nil {
		if *maxArgs < minArgs {
			return fmt.Errorf("maxArgs (%d) must be >= minArgs (%d)", *maxArgs, minArgs)
		}
	}

	// Check for duplicate names
	if _, exists := r.functions[name]; exists {
		return fmt.Errorf("function %q is already registered", name)
	}

	// Validate handler signature
	if err := r.validateHandlerSignature(handler, maxArgs); err != nil {
		return fmt.Errorf("invalid handler signature for function %q: %w", name, err)
	}

	r.functions[name] = &TemplateFunction{
		Name:      name,
		MinArgs:   minArgs,
		MaxArgs:   maxArgs,
		Validator: validator,
		Handler:   handler,
	}

	return nil
}

// Get retrieves a template function by name
func (r *FunctionRegistry) Get(name string) (*TemplateFunction, bool) {
	fn, ok := r.functions[name]
	return fn, ok
}

// validateHandlerSignature checks if the handler function has the correct signature
func (r *FunctionRegistry) validateHandlerSignature(handler interface{}, maxArgs *int) error {
	handlerType := reflect.TypeOf(handler)
	
	// Must be a function
	if handlerType.Kind() != reflect.Func {
		return fmt.Errorf("handler must be a function, got %s", handlerType.Kind())
	}

	// Check parameter count and types
	if maxArgs != nil {
		// Fixed args: func(*Context, string, string, ...)
		expectedParams := *maxArgs + 1 // +1 for Context
		if handlerType.NumIn() != expectedParams {
			return fmt.Errorf("handler must have %d parameters (*Context + %d string args), got %d", expectedParams, *maxArgs, handlerType.NumIn())
		}

		// First param must be *Context
		if handlerType.NumIn() > 0 {
			firstParam := handlerType.In(0)
			contextType := reflect.TypeOf((*Context)(nil))
			if firstParam != contextType {
				return fmt.Errorf("first parameter must be *Context, got %s", firstParam)
			}
		}

		// Remaining params must be string
		for i := 1; i < handlerType.NumIn(); i++ {
			paramType := handlerType.In(i)
			if paramType.Kind() != reflect.String {
				return fmt.Errorf("parameter %d must be string, got %s", i, paramType.Kind())
			}
		}
	} else {
		// Variadic: func(*Context, []string)
		if handlerType.NumIn() != 2 {
			return fmt.Errorf("variadic handler must have 2 parameters (*Context, []string), got %d", handlerType.NumIn())
		}

		// First param must be *Context
		firstParam := handlerType.In(0)
		contextType := reflect.TypeOf((*Context)(nil))
		if firstParam != contextType {
			return fmt.Errorf("first parameter must be *Context, got %s", firstParam)
		}

		// Second param must be []string
		secondParam := handlerType.In(1)
		if secondParam.Kind() != reflect.Slice || secondParam.Elem().Kind() != reflect.String {
			return fmt.Errorf("second parameter must be []string, got %s", secondParam)
		}
	}

	return nil
}

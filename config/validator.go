package config

import (
	"fmt"

	"github.com/cassdeckard/tviewyaml/keys"
)

// Validator validates configuration structures
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateApp validates an application configuration
func (v *Validator) ValidateApp(config *AppConfig) error {
	// Root element validation (currently only supports "pages" type)
	if config.Application.Root.Type != "pages" {
		return fmt.Errorf("application root type must be 'pages', got: %s", config.Application.Root.Type)
	}

	if len(config.Application.Root.Pages) == 0 {
		return fmt.Errorf("application root must contain at least one page")
	}

	// Validate page references
	for i, page := range config.Application.Root.Pages {
		if page.Name == "" {
			return fmt.Errorf("page %d is missing name", i)
		}
		if page.Ref == "" {
			return fmt.Errorf("page %s is missing ref", page.Name)
		}
	}

	// Validate key bindings
	for i, binding := range config.Application.GlobalKeyBindings {
		if binding.Key == "" {
			return fmt.Errorf("key binding %d is missing key", i)
		}
		if _, _, _, err := keys.ParseKey(binding.Key); err != nil {
			return fmt.Errorf("key binding %d has invalid key %q: %w", i, binding.Key, err)
		}
		if binding.Action == "" {
			return fmt.Errorf("key binding %d is missing action", i)
		}
	}

	return nil
}

// ValidatePage validates a page configuration
func (v *Validator) ValidatePage(config *PageConfig) error {
	if config.Type == "" {
		return fmt.Errorf("page type is required")
	}

	// Validate based on type
	switch config.Type {
	case "list":
		if len(config.ListItems) == 0 && len(config.Items) == 0 {
			return fmt.Errorf("list type requires listItems or items")
		}
	case "flex":
		if len(config.Items) == 0 {
			return fmt.Errorf("flex type requires items")
		}
	case "form":
		if len(config.FormItems) == 0 {
			return fmt.Errorf("form type requires formItems")
		}
	case "table":
		if config.TableData == nil {
			return fmt.Errorf("table type requires tableData")
		}
	}

	return nil
}

// ValidatePrimitive validates a primitive configuration
func (v *Validator) ValidatePrimitive(prim *Primitive) error {
	if prim.Type == "" {
		return fmt.Errorf("primitive type is required")
	}

	// Add more validation as needed
	return nil
}

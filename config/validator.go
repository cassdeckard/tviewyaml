package config

import (
	"fmt"
)

// Validator validates configuration structures
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateRoot validates a root configuration
func (v *Validator) ValidateRoot(config *RootConfig) error {
	if config.Root.Type != "pages" {
		return fmt.Errorf("root type must be 'pages', got: %s", config.Root.Type)
	}

	if len(config.Root.Pages) == 0 {
		return fmt.Errorf("root must contain at least one page")
	}

	for i, page := range config.Root.Pages {
		if page.Name == "" {
			return fmt.Errorf("page %d is missing name", i)
		}
		if page.Ref == "" {
			return fmt.Errorf("page %s is missing ref", page.Name)
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

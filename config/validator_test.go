package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateApp(t *testing.T) {
	tests := []struct {
		name    string
		config  *AppConfig
		wantErr bool
		errContains string
	}{
		// Valid cases
		{
			name: "valid config",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "pages",
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid with key bindings",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "pages",
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
						},
					},
					GlobalKeyBindings: []KeyBinding{
						{Key: "Escape", Action: "{{ stopApp }}"},
						{Key: "Ctrl+Q", Action: "{{ stopApp }}"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid with multiple pages",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "pages",
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
							{Name: "form", Ref: "form.yaml"},
						},
					},
				},
			},
			wantErr: false,
		},

		// Invalid root type
		{
			name: "invalid root type",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "invalid",
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
						},
					},
				},
			},
			wantErr: true,
			errContains: "root type must be 'pages'",
		},

		// Empty pages
		{
			name: "no pages",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type:  "pages",
						Pages: []PageRef{},
					},
				},
			},
			wantErr: true,
			errContains: "must contain at least one page",
		},

		// Page reference validation
		{
			name: "page missing name",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "pages",
						Pages: []PageRef{
							{Name: "", Ref: "main.yaml"},
						},
					},
				},
			},
			wantErr: true,
			errContains: "page 0 is missing name",
		},
		{
			name: "page missing ref",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "pages",
						Pages: []PageRef{
							{Name: "main", Ref: ""},
						},
					},
				},
			},
			wantErr: true,
			errContains: "page main is missing ref",
		},
		{
			name: "multiple pages with missing refs",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "pages",
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
							{Name: "form", Ref: ""},
						},
					},
				},
			},
			wantErr: true,
			errContains: "page form is missing ref",
		},

		// Key binding validation
		{
			name: "key binding missing key",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "pages",
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
						},
					},
					GlobalKeyBindings: []KeyBinding{
						{Key: "", Action: "{{ stopApp }}"},
					},
				},
			},
			wantErr: true,
			errContains: "key binding 0 is missing key",
		},
		{
			name: "key binding missing action",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "pages",
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
						},
					},
					GlobalKeyBindings: []KeyBinding{
						{Key: "Escape", Action: ""},
					},
				},
			},
			wantErr: true,
			errContains: "key binding 0 is missing action",
		},
		{
			name: "key binding invalid key",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "pages",
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
						},
					},
					GlobalKeyBindings: []KeyBinding{
						{Key: "InvalidKey", Action: "{{ stopApp }}"},
					},
				},
			},
			wantErr: true,
			errContains: "key binding 0 has invalid key",
		},
		{
			name: "multiple key bindings with errors",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Type: "pages",
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
						},
					},
					GlobalKeyBindings: []KeyBinding{
						{Key: "Escape", Action: "{{ stopApp }}"},
						{Key: "InvalidKey", Action: "{{ stopApp }}"},
					},
				},
			},
			wantErr: true,
			errContains: "key binding 1 has invalid key",
		},
	}

	validator := NewValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateApp(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("ValidateApp() error = %v, want error containing %q", err, tt.errContains)
					}
				}
			}
		})
	}
}

func TestValidatePage(t *testing.T) {
	tests := []struct {
		name        string
		config      *PageConfig
		wantErr     bool
		errContains string
	}{
		// Valid cases
		{
			name: "valid list with listItems",
			config: &PageConfig{
				Type: "list",
				ListItems: []ListItem{
					{MainText: "Item 1"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid list with items",
			config: &PageConfig{
				Type: "list",
				Items: []FlexItem{
					{Primitive: &Primitive{Type: "textView"}},
				},
			},
			wantErr: false,
		},
		{
			name: "valid flex",
			config: &PageConfig{
				Type: "flex",
				Items: []FlexItem{
					{Primitive: &Primitive{Type: "textView"}},
				},
			},
			wantErr: false,
		},
		{
			name: "valid form",
			config: &PageConfig{
				Type: "form",
				FormItems: []FormItem{
					{Type: "inputfield", Label: "Name"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid table",
			config: &PageConfig{
				Type: "table",
				TableData: &TableData{
					Headers: []string{"Col1"},
					Rows:    [][]string{{"Data"}},
				},
			},
			wantErr: false,
		},
		{
			name: "valid treeView",
			config: &PageConfig{
				Type: "treeView",
				Nodes: []TreeNode{}, // Empty tree is valid
			},
			wantErr: false,
		},
		{
			name: "valid treeView with nodes",
			config: &PageConfig{
				Type: "treeView",
				Nodes: []TreeNode{
					{Name: "root", Text: "Root"},
				},
			},
			wantErr: false,
		},

		// Missing type
		{
			name: "missing type",
			config: &PageConfig{
				Type: "",
			},
			wantErr: true,
			errContains: "page type is required",
		},

		// List validation
		{
			name: "list without listItems or items",
			config: &PageConfig{
				Type:      "list",
				ListItems: []ListItem{},
				Items:     []FlexItem{},
			},
			wantErr: true,
			errContains: "list type requires listItems or items",
		},

		// Flex validation
		{
			name: "flex without items",
			config: &PageConfig{
				Type:  "flex",
				Items: []FlexItem{},
			},
			wantErr: true,
			errContains: "flex type requires items",
		},

		// Form validation
		{
			name: "form without formItems",
			config: &PageConfig{
				Type:      "form",
				FormItems: []FormItem{},
			},
			wantErr: true,
			errContains: "form type requires formItems",
		},

		// Table validation
		{
			name: "table without tableData",
			config: &PageConfig{
				Type:      "table",
				TableData: nil,
			},
			wantErr: true,
			errContains: "table type requires tableData",
		},
	}

	validator := NewValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePage(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("ValidatePage() error = %v, want error containing %q", err, tt.errContains)
					}
				}
			}
		})
	}
}

func TestValidateAppRefs(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create test files
	mainYAML := filepath.Join(tmpDir, "main.yaml")
	if err := os.WriteFile(mainYAML, []byte("type: list\nlistItems: []"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	formYAML := filepath.Join(tmpDir, "form.yaml")
	if err := os.WriteFile(formYAML, []byte("type: form\nformItems: []"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		config      *AppConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "all refs exist",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
							{Name: "form", Ref: "form.yaml"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing ref file",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Pages: []PageRef{
							{Name: "main", Ref: "main.yaml"},
							{Name: "missing", Ref: "missing.yaml"},
						},
					},
				},
			},
			wantErr: true,
			errContains: "file does not exist",
		},
		{
			name: "all refs missing",
			config: &AppConfig{
				Application: ApplicationElement{
					Root: RootElement{
						Pages: []PageRef{
							{Name: "missing1", Ref: "missing1.yaml"},
							{Name: "missing2", Ref: "missing2.yaml"},
						},
					},
				},
			},
			wantErr: true,
			errContains: "file does not exist",
		},
	}

	validator := NewValidator()
	loader := NewLoader(tmpDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAppRefs(tt.config, loader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAppRefs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("ValidateAppRefs() error = %v, want error containing %q", err, tt.errContains)
					}
				}
			}
		})
	}
}

func TestValidatePrimitive(t *testing.T) {
	tests := []struct {
		name        string
		prim        *Primitive
		wantErr     bool
		errContains string
	}{
		{
			name: "valid primitive",
			prim: &Primitive{
				Type: "textView",
			},
			wantErr: false,
		},
		{
			name: "missing type",
			prim: &Primitive{
				Type: "",
			},
			wantErr: true,
			errContains: "primitive type is required",
		},
	}

	validator := NewValidator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePrimitive(tt.prim)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePrimitive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("ValidatePrimitive() error = %v, want error containing %q", err, tt.errContains)
					}
				}
			}
		})
	}
}

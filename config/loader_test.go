package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadApp(t *testing.T) {
	tmpDir := t.TempDir()

	validYAML := `application:
  root:
    type: pages
    pages:
      - name: main
        ref: main.yaml
`

	invalidYAML := `application:
  root:
    type: pages
    pages:
      - name: main
        ref: main.yaml
    invalid: [unclosed bracket
`

	tests := []struct {
		name        string
		setup       func() string // returns filename
		wantErr     bool
		errContains string
		validate    func(*AppConfig) bool // optional validation
	}{
		{
			name: "valid app config",
			setup: func() string {
				filename := filepath.Join(tmpDir, "app.yaml")
				if err := os.WriteFile(filename, []byte(validYAML), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return "app.yaml"
			},
			wantErr: false,
			validate: func(cfg *AppConfig) bool {
				return cfg.Application.Root.Type == "pages" &&
					len(cfg.Application.Root.Pages) == 1 &&
					cfg.Application.Root.Pages[0].Name == "main"
			},
		},
		{
			name: "file not found",
			setup: func() string {
				return "nonexistent.yaml"
			},
			wantErr: true,
			errContains: "failed to read app config",
		},
		{
			name: "invalid YAML",
			setup: func() string {
				filename := filepath.Join(tmpDir, "invalid.yaml")
				if err := os.WriteFile(filename, []byte(invalidYAML), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return "invalid.yaml"
			},
			wantErr: true,
			errContains: "failed to parse app config",
		},
		{
			name: "empty file",
			setup: func() string {
				filename := filepath.Join(tmpDir, "empty.yaml")
				if err := os.WriteFile(filename, []byte(""), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return "empty.yaml"
			},
			wantErr: false, // Empty YAML is valid (results in zero values)
		},
	}

	loader := NewLoader(tmpDir)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.setup()
			cfg, err := loader.LoadApp(filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadApp(%q) error = %v, wantErr %v", filename, err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("LoadApp(%q) error = %v, want error containing %q", filename, err, tt.errContains)
					}
				}
				return
			}
			if tt.validate != nil && !tt.validate(cfg) {
				t.Errorf("LoadApp(%q) returned invalid config", filename)
			}
		})
	}
}

func TestLoadPage(t *testing.T) {
	tmpDir := t.TempDir()

	validYAML := `type: list
listItems:
  - mainText: Item 1
`

	invalidYAML := `type: list
listItems:
  - mainText: [unclosed bracket
`

	tests := []struct {
		name        string
		setup       func() string // returns ref
		wantErr     bool
		errContains string
		validate    func(*PageConfig) bool // optional validation
	}{
		{
			name: "valid page config",
			setup: func() string {
				ref := "page.yaml"
				filename := filepath.Join(tmpDir, ref)
				if err := os.WriteFile(filename, []byte(validYAML), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return ref
			},
			wantErr: false,
			validate: func(cfg *PageConfig) bool {
				return cfg.Type == "list" &&
					len(cfg.ListItems) == 1 &&
					cfg.ListItems[0].MainText == "Item 1"
			},
		},
		{
			name: "file not found",
			setup: func() string {
				return "nonexistent.yaml"
			},
			wantErr: true,
			errContains: "failed to read page config",
		},
		{
			name: "invalid YAML",
			setup: func() string {
				ref := "invalid.yaml"
				filename := filepath.Join(tmpDir, ref)
				if err := os.WriteFile(filename, []byte(invalidYAML), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return ref
			},
			wantErr: true,
			errContains: "failed to parse page config",
		},
		{
			name: "empty file",
			setup: func() string {
				ref := "empty.yaml"
				filename := filepath.Join(tmpDir, ref)
				if err := os.WriteFile(filename, []byte(""), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return ref
			},
			wantErr: false, // Empty YAML is valid (results in zero values)
		},
		{
			name: "nested path",
			setup: func() string {
				subDir := filepath.Join(tmpDir, "subdir")
				if err := os.MkdirAll(subDir, 0755); err != nil {
					t.Fatalf("Failed to create subdirectory: %v", err)
				}
				ref := "subdir/page.yaml"
				filename := filepath.Join(tmpDir, ref)
				if err := os.WriteFile(filename, []byte(validYAML), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return ref
			},
			wantErr: false,
			validate: func(cfg *PageConfig) bool {
				return cfg.Type == "list"
			},
		},
	}

	loader := NewLoader(tmpDir)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := tt.setup()
			cfg, err := loader.LoadPage(ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPage(%q) error = %v, wantErr %v", ref, err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("LoadPage(%q) error = %v, want error containing %q", ref, err, tt.errContains)
					}
				}
				return
			}
			if tt.validate != nil && !tt.validate(cfg) {
				t.Errorf("LoadPage(%q) returned invalid config", ref)
			}
		})
	}
}

func TestLoadPageDirect(t *testing.T) {
	tmpDir := t.TempDir()

	validYAML := `type: form
formItems:
  - type: inputfield
    label: Name
`

	tests := []struct {
		name        string
		setup       func() string // returns absolute path
		wantErr     bool
		errContains string
		validate    func(*PageConfig) bool // optional validation
	}{
		{
			name: "valid page config with absolute path",
			setup: func() string {
				filename := filepath.Join(tmpDir, "page.yaml")
				if err := os.WriteFile(filename, []byte(validYAML), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				absPath, err := filepath.Abs(filename)
				if err != nil {
					t.Fatalf("Failed to get absolute path: %v", err)
				}
				return absPath
			},
			wantErr: false,
			validate: func(cfg *PageConfig) bool {
				return cfg.Type == "form" &&
					len(cfg.FormItems) == 1 &&
					cfg.FormItems[0].Label == "Name"
			},
		},
		{
			name: "file not found",
			setup: func() string {
				return filepath.Join(tmpDir, "nonexistent.yaml")
			},
			wantErr: true,
			errContains: "failed to read page config",
		},
		{
			name: "relative path",
			setup: func() string {
				filename := filepath.Join(tmpDir, "page.yaml")
				if err := os.WriteFile(filename, []byte(validYAML), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				// Use relative path from current directory
				relPath, err := filepath.Rel(".", filename)
				if err != nil {
					// If relative path fails, use absolute
					relPath, _ = filepath.Abs(filename)
				}
				return relPath
			},
			wantErr: false,
			validate: func(cfg *PageConfig) bool {
				return cfg.Type == "form"
			},
		},
	}

	loader := NewLoader(tmpDir) // basePath doesn't matter for LoadPageDirect
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			cfg, err := loader.LoadPageDirect(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPageDirect(%q) error = %v, wantErr %v", path, err, tt.wantErr)
				return
			}
			if err != nil {
				if tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("LoadPageDirect(%q) error = %v, want error containing %q", path, err, tt.errContains)
					}
				}
				return
			}
			if tt.validate != nil && !tt.validate(cfg) {
				t.Errorf("LoadPageDirect(%q) returned invalid config", path)
			}
		})
	}
}

func TestRefExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	existingFile := filepath.Join(tmpDir, "existing.yaml")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	nestedFile := filepath.Join(subDir, "nested.yaml")
	if err := os.WriteFile(nestedFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create nested test file: %v", err)
	}

	tests := []struct {
		name     string
		ref      string
		want     bool
		desc     string
	}{
		{
			name: "existing file",
			ref:  "existing.yaml",
			want: true,
			desc: "file exists in base path",
		},
		{
			name: "nonexistent file",
			ref:  "nonexistent.yaml",
			want: false,
			desc: "file does not exist",
		},
		{
			name: "nested path",
			ref:  "subdir/nested.yaml",
			want: true,
			desc: "file exists in nested directory",
		},
		{
			name: "nested nonexistent",
			ref:  "subdir/missing.yaml",
			want: false,
			desc: "file does not exist in nested directory",
		},
		{
			name: "empty ref",
			ref:  "",
			want: true, // filepath.Join(basePath, "") returns basePath, which exists
			desc: "empty ref string resolves to base path",
		},
		{
			name: "directory instead of file",
			ref:  "subdir",
			want: true, // os.Stat returns no error for directories
			desc: "ref points to directory",
		},
	}

	loader := NewLoader(tmpDir)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := loader.RefExists(tt.ref)
			if got != tt.want {
				t.Errorf("RefExists(%q) = %v, want %v (%s)", tt.ref, got, tt.want, tt.desc)
			}
		})
	}
}

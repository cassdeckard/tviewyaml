package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader handles loading YAML configuration files
type Loader struct {
	basePath string
}

// NewLoader creates a new config loader with a base path
func NewLoader(basePath string) *Loader {
	return &Loader{basePath: basePath}
}

// LoadRoot loads the root configuration file
func (l *Loader) LoadRoot(filename string) (*RootConfig, error) {
	path := filepath.Join(l.basePath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read root config %s: %w", path, err)
	}

	var config RootConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse root config %s: %w", path, err)
	}

	return &config, nil
}

// LoadPage loads a page configuration file
func (l *Loader) LoadPage(ref string) (*PageConfig, error) {
	path := filepath.Join(l.basePath, ref)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read page config %s: %w", path, err)
	}

	var config PageConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse page config %s: %w", path, err)
	}

	return &config, nil
}

// LoadPageDirect loads a page config from an absolute or relative path
func (l *Loader) LoadPageDirect(path string) (*PageConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read page config %s: %w", path, err)
	}

	var config PageConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse page config %s: %w", path, err)
	}

	return &config, nil
}

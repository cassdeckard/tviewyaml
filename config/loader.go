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

// LoadApp loads the application configuration file
func (l *Loader) LoadApp(filename string) (*AppConfig, error) {
	path := filepath.Join(l.basePath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read app config %s: %w", path, err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse app config %s: %w", path, err)
	}

	return &config, nil
}

// RefExists returns true if the page ref file exists under the loader's base path.
func (l *Loader) RefExists(ref string) bool {
	path := filepath.Join(l.basePath, ref)
	_, err := os.Stat(path)
	return err == nil
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

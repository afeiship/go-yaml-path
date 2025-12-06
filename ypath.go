package ypath

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// YPath represents a YAML document with path navigation capabilities
type YPath struct {
	data interface{}
}

// New creates a new YPath instance from YAML bytes
func New(yamlData []byte) (*YPath, error) {
	var data interface{}
	err := yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return &YPath{data: data}, nil
}

// NewFromString creates a new YPath instance from a YAML string
func NewFromString(yamlStr string) (*YPath, error) {
	return New([]byte(yamlStr))
}

// NewFromFile creates a new YPath instance from a YAML file
func NewFromFile(filename string) (*YPath, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return New(data)
}
package ypath

import (
	"fmt"
	"os"

	"github.com/afeiship/go-dotpath"
	"gopkg.in/yaml.v3"
)

// YPath represents a YAML document with path navigation capabilities
// Internally uses go-dotpath for implementation
type YPath struct {
	dp *dotpath.DotPath
}

// New creates a new YPath instance from YAML bytes
func New(yamlData []byte) (*YPath, error) {
	var data any
	err := yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Handle empty YAML or non-map YAML
	var mapData map[string]any
	if data == nil {
		mapData = make(map[string]any)
	} else if m, ok := data.(map[string]any); ok {
		mapData = m
	} else {
		// If YAML doesn't parse to a map, wrap it
		mapData = map[string]any{"data": data}
	}

	return &YPath{dp: dotpath.New(mapData)}, nil
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

// NewFromMap creates a new YPath instance from a map interface
func NewFromMap(data map[string]any) *YPath {
	return &YPath{dp: dotpath.New(data)}
}

// Data returns the underlying map data
func (yp *YPath) Data() map[string]any {
	return yp.dp.Data()
}
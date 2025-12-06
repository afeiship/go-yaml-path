package ypath

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

// Get retrieves a value by dot notation path (e.g., "server.host", "database.connection.pool.max")
func (yp *YPath) Get(path string) interface{} {
	if path == "" || path == "." {
		return yp.data
	}

	parts := strings.Split(path, ".")
	current := yp.data

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil
			}
		case map[interface{}]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil
			}
		case []interface{}:
			// Handle array access with numeric index
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= len(v) {
				return nil
			}
			current = v[index]
		default:
			return nil
		}
	}

	return current
}

// GetString retrieves a string value by path
func (yp *YPath) GetString(path string) string {
	val := yp.Get(path)
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetInt retrieves an integer value by path
func (yp *YPath) GetInt(path string) int {
	val := yp.Get(path)
	if val == nil {
		return 0
	}

	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// GetBool retrieves a boolean value by path
func (yp *YPath) GetBool(path string) bool {
	val := yp.Get(path)
	if val == nil {
		return false
	}

	switch v := val.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true" || v == "1"
	case int:
		return v != 0
	case float64:
		return v != 0
	}
	return false
}

// GetFloat64 retrieves a float64 value by path
func (yp *YPath) GetFloat64(path string) float64 {
	val := yp.Get(path)
	if val == nil {
		return 0.0
	}

	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.0
}

// Exists checks if a path exists in the YAML data
func (yp *YPath) Exists(path string) bool {
	return yp.Get(path) != nil
}

// GetAll retrieves all values matching a wildcard path (e.g., "servers.*.host")
func (yp *YPath) GetAll(path string) []interface{} {
	if !strings.Contains(path, "*") {
		if val := yp.Get(path); val != nil {
			return []interface{}{val}
		}
		return nil
	}

	parts := strings.Split(path, ".")
	var results []interface{}
	yp.collectWildcardValues(yp.data, parts, 0, &results)
	return results
}

// collectWildcardValues recursively collects values from wildcard paths
func (yp *YPath) collectWildcardValues(current interface{}, parts []string, index int, results *[]interface{}) {
	if index >= len(parts) {
		*results = append(*results, current)
		return
	}

	part := parts[index]

	if part == "*" {
		switch v := current.(type) {
		case map[string]interface{}:
			for _, val := range v {
				yp.collectWildcardValues(val, parts, index+1, results)
			}
		case map[interface{}]interface{}:
			for _, val := range v {
				yp.collectWildcardValues(val, parts, index+1, results)
			}
		case []interface{}:
			for _, val := range v {
				yp.collectWildcardValues(val, parts, index+1, results)
			}
		}
		return
	}

	switch v := current.(type) {
	case map[string]interface{}:
		if next, ok := v[part]; ok {
			yp.collectWildcardValues(next, parts, index+1, results)
		}
	case map[interface{}]interface{}:
		if next, ok := v[part]; ok {
			yp.collectWildcardValues(next, parts, index+1, results)
		}
	case []interface{}:
		if index == len(parts)-1 && part == "*" {
			*results = append(*results, v...)
		} else {
			// Handle array access with numeric index
			if arrayIndex, err := strconv.Atoi(part); err == nil && arrayIndex >= 0 && arrayIndex < len(v) {
				yp.collectWildcardValues(v[arrayIndex], parts, index+1, results)
			}
		}
	}
}

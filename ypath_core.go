package ypath

import (
	"fmt"
	"strconv"
	"strings"
)

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
package ypath

import (
	"fmt"
	"strconv"
	"strings"
)

// Get retrieves a value by dot notation path (e.g., "server.host", "database.connection.pool.max")
// If no path is provided or path is empty/dot, returns the entire YAML data structure
func (yp *YPath) Get(path ...string) any {
	if len(path) == 0 || path[0] == "" || path[0] == "." {
		return yp.dp.Data()
	}

	p := path[0]

	// Try dotpath first - it might handle some array cases
	if val, exists := yp.dp.Get(p); exists {
		return val
	}

	// If dotpath didn't find it, try manual array access
	if yp.containsArrayAccess(p) {
		return yp.getArrayValue(p)
	}

	return nil
}

// containsArrayAccess checks if a path contains numeric array indices
func (yp *YPath) containsArrayAccess(path string) bool {
	parts := strings.Split(path, ".")
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err == nil {
			return true
		}
	}
	return false
}

// getArrayValue handles array access using manual navigation
func (yp *YPath) getArrayValue(path string) any {
	parts := strings.Split(path, ".")
	var current any = yp.dp.Data()

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]any:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil
			}
		case map[any]any:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil
			}
		case []any:
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
	if str, ok := val.(string); ok {
		return str
	}
	// Convert non-string values to string for compatibility
	return fmt.Sprintf("%v", val)
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
		return v == "true" || v == "1"
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
	// Always use our Get method for consistency, including for "." and ""
	return yp.Get(path) != nil
}
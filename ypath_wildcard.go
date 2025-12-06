package ypath

import (
	"strconv"
	"strings"
)

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
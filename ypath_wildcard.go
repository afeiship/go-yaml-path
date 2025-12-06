package ypath

import (
	"strconv"
	"strings"
)

// GetAll retrieves all values matching a wildcard path (e.g., "servers.*.host")
func (yp *YPath) GetAll(path string) []any {
	if !strings.Contains(path, "*") {
		// Non-wildcard paths: use our optimized Get method
		if val := yp.Get(path); val != nil {
			return []any{val}
		}
		return nil
	}

	// Wildcard paths: use custom logic with dotpath data
	parts := strings.Split(path, ".")
	var results []any
	yp.collectWildcardValues(yp.dp.Data(), parts, 0, &results)
	return results
}

// collectWildcardValues recursively collects values from wildcard paths
func (yp *YPath) collectWildcardValues(current any, parts []string, index int, results *[]any) {
	if index >= len(parts) {
		*results = append(*results, current)
		return
	}

	part := parts[index]

	if part == "*" {
		switch v := current.(type) {
		case map[string]any:
			for _, val := range v {
				yp.collectWildcardValues(val, parts, index+1, results)
			}
		case map[any]any:
			for _, val := range v {
				yp.collectWildcardValues(val, parts, index+1, results)
			}
		case []any:
			for _, val := range v {
				yp.collectWildcardValues(val, parts, index+1, results)
			}
		}
		return
	}

	switch v := current.(type) {
	case map[string]any:
		if next, ok := v[part]; ok {
			yp.collectWildcardValues(next, parts, index+1, results)
		}
	case map[any]any:
		if next, ok := v[part]; ok {
			yp.collectWildcardValues(next, parts, index+1, results)
		}
	case []any:
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
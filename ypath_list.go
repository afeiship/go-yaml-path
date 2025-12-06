package ypath

import (
	"fmt"
	"strconv"
	"strings"
)

// GetStringList retrieves a string list by path, converts other types to strings
func (yp *YPath) GetStringList(path string) []string {
	var val interface{}

	// Use GetAll for wildcard paths
	if strings.Contains(path, "*") {
		val = yp.GetAll(path)
	} else {
		val = yp.Get(path)
	}

	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result
	case []string:
		return v
	case string:
		return []string{v}
	default:
		return []string{fmt.Sprintf("%v", v)}
	}
}

// GetIntList retrieves an integer list by path, with type conversion
func (yp *YPath) GetIntList(path string) []int {
	var val interface{}

	// Use GetAll for wildcard paths
	if strings.Contains(path, "*") {
		val = yp.GetAll(path)
	} else {
		val = yp.Get(path)
	}

	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case []interface{}:
		result := make([]int, 0, len(v))
		for _, item := range v {
			switch num := item.(type) {
			case int:
				result = append(result, num)
			case int64:
				result = append(result, int(num))
			case float64:
				result = append(result, int(num))
			case string:
				if i, err := strconv.Atoi(num); err == nil {
					result = append(result, i)
				}
			}
		}
		return result
	case []int:
		return v
	case []int64:
		result := make([]int, len(v))
		for i, num := range v {
			result[i] = int(num)
		}
		return result
	case int:
		return []int{v}
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return []int{i}
		}
	}
	return nil
}

// GetBoolList retrieves a boolean list by path, with type conversion
func (yp *YPath) GetBoolList(path string) []bool {
	var val interface{}

	// Use GetAll for wildcard paths
	if strings.Contains(path, "*") {
		val = yp.GetAll(path)
	} else {
		val = yp.Get(path)
	}

	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case []interface{}:
		result := make([]bool, 0, len(v))
		for _, item := range v {
			switch b := item.(type) {
			case bool:
				result = append(result, b)
			case string:
				result = append(result, strings.ToLower(b) == "true" || b == "1")
			case int:
				result = append(result, b != 0)
			case float64:
				result = append(result, b != 0)
			}
		}
		return result
	case []bool:
		return v
	case bool:
		return []bool{v}
	case string:
		return []bool{strings.ToLower(v) == "true" || v == "1"}
	}
	return nil
}

// GetFloat64List retrieves a float64 list by path, with type conversion
func (yp *YPath) GetFloat64List(path string) []float64 {
	var val interface{}

	// Use GetAll for wildcard paths
	if strings.Contains(path, "*") {
		val = yp.GetAll(path)
	} else {
		val = yp.Get(path)
	}

	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case []interface{}:
		result := make([]float64, 0, len(v))
		for _, item := range v {
			switch num := item.(type) {
			case float64:
				result = append(result, num)
			case int:
				result = append(result, float64(num))
			case int64:
				result = append(result, float64(num))
			case string:
				if f, err := strconv.ParseFloat(num, 64); err == nil {
					result = append(result, f)
				}
			}
		}
		return result
	case []float64:
		return v
	case []int:
		result := make([]float64, len(v))
		for i, num := range v {
			result[i] = float64(num)
		}
		return result
	case float64:
		return []float64{v}
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return []float64{f}
		}
	}
	return nil
}

// GetList is a generic method to get a list of any type
func (yp *YPath) GetList(path string) []interface{} {
	var val interface{}

	// Use GetAll for wildcard paths
	if strings.Contains(path, "*") {
		val = yp.GetAll(path)
	} else {
		val = yp.Get(path)
	}

	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case []interface{}:
		return v
	case []string:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result
	case []int:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result
	case []bool:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result
	case []float64:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result
	default:
		// Convert single value to list
		return []interface{}{v}
	}
}
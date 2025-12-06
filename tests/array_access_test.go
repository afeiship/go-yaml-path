package tests

import "testing"

func TestArrayAccess(t *testing.T) {
	yp := createSampleYPath(t)

	tests := []struct {
		path     string
		expected string
	}{
		{"features.0", "authentication"},
		{"features.1", "logging"},
		{"features.2", "monitoring"},
		{"servers.0.host", "server1.example.com"},
		{"servers.1.port", "8002"},
	}

	for _, test := range tests {
		result := yp.GetString(test.path)
		if result != test.expected {
			t.Errorf("Array access %s = %q, expected %q", test.path, result, test.expected)
		}
	}
}

func TestArrayIndexOutOfBounds(t *testing.T) {
	yp := createSampleYPath(t)

	// Test out of bounds access
	result := yp.GetString("features.10")
	if result != "" {
		t.Errorf("Expected empty string for out of bounds access, got %q", result)
	}

	resultInt := yp.GetInt("servers.10.port")
	if resultInt != 0 {
		t.Errorf("Expected 0 for out of bounds int access, got %d", resultInt)
	}
}

func TestEmptyArrayOperations(t *testing.T) {
	emptyYAML := `
empty_array: []
`

	yp := createYPath(t, emptyYAML)

	// Test empty array access
	result := yp.GetString("empty_array.0")
	if result != "" {
		t.Errorf("Expected empty string for empty array access, got %q", result)
	}

	// Test array length check
	array := yp.Get("empty_array")
	if array == nil {
		t.Error("Should get empty array, not nil")
	}
}
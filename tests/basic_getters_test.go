package tests

import "testing"

func TestGetString(t *testing.T) {
	yp := createSampleYPath(t)

	tests := []struct {
		path     string
		expected string
	}{
		{"server.host", "localhost"},
		{"database.credentials.username", "admin"},
		{"server.port", "8080"},
		{"nonexistent.path", ""},
	}

	for _, test := range tests {
		result := yp.GetString(test.path)
		if result != test.expected {
			t.Errorf("GetString(%s) = %q, expected %q", test.path, result, test.expected)
		}
	}
}

func TestGetInt(t *testing.T) {
	yp := createSampleYPath(t)

	tests := []struct {
		path     string
		expected int
	}{
		{"server.port", 8080},
		{"database.connection.pool.max", 10},
		{"numbers.int_val", 42},
		{"numbers.float_val", 3},
		{"numbers.string_num", 123},
		{"nonexistent.path", 0},
	}

	for _, test := range tests {
		result := yp.GetInt(test.path)
		if result != test.expected {
			t.Errorf("GetInt(%s) = %d, expected %d", test.path, result, test.expected)
		}
	}
}

func TestGetBool(t *testing.T) {
	yp := createSampleYPath(t)

	tests := []struct {
		path     string
		expected bool
	}{
		{"server.ssl", true},
		{"numbers.bool_true", true},
		{"numbers.bool_one", true},
		{"numbers.bool_false", false},
		{"servers.0.active", true},
		{"servers.1.active", false},
		{"nonexistent.path", false},
	}

	for _, test := range tests {
		result := yp.GetBool(test.path)
		if result != test.expected {
			t.Errorf("GetBool(%s) = %t, expected %t", test.path, result, test.expected)
		}
	}
}

func TestGetFloat64(t *testing.T) {
	yp := createSampleYPath(t)

	tests := []struct {
		path     string
		expected float64
	}{
		{"database.connection.timeout", 30.5},
		{"numbers.float_val", 3.14},
		{"numbers.int_val", 42.0},
		{"nonexistent.path", 0.0},
	}

	for _, test := range tests {
		result := yp.GetFloat64(test.path)
		if result != test.expected {
			t.Errorf("GetFloat64(%s) = %f, expected %f", test.path, result, test.expected)
		}
	}
}

func TestGet(t *testing.T) {
	yp := createSampleYPath(t)

	// Test getting different types
	serverHost := yp.Get("server.host")
	if serverHost != "localhost" {
		t.Errorf("Get(server.host) = %v, expected localhost", serverHost)
	}

	serverPort := yp.Get("server.port")
	if serverPort != 8080 {
		t.Errorf("Get(server.port) = %v, expected 8080", serverPort)
	}

	serverSSL := yp.Get("server.ssl")
	if serverSSL != true {
		t.Errorf("Get(server.ssl) = %v, expected true", serverSSL)
	}

	// Test non-existent path
	nonExistent := yp.Get("nonexistent.path")
	if nonExistent != nil {
		t.Errorf("Get(nonexistent.path) = %v, expected nil", nonExistent)
	}

	// Test root access
	root := yp.Get("")
	if root == nil {
		t.Error("Get(\"\") should return root data")
	}

	rootDot := yp.Get(".")
	if rootDot == nil {
		t.Error("Get(\".\") should return root data")
	}
}

func TestGetWithoutParameters(t *testing.T) {
	yp := createSampleYPath(t)

	// Test Get() with no parameters returns root data
	data := yp.Get()
	if data == nil {
		t.Error("Get() should return root data")
	}

	// Test type is correct
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		t.Errorf("Get() should return map[string]interface{}, got %T", data)
	}

	// Verify Get() returns the same as Get("") - check size equality
	rootFromGet := yp.Get("")
	rootMap, rootOk := rootFromGet.(map[string]interface{})

	if !rootOk {
		t.Error("Get(\"\") should return map[string]interface{}")
	} else if len(dataMap) != len(rootMap) {
		t.Error("Get() and Get(\"\") should return maps of same size")
	}

	// Verify Get() returns the same as Get(".")
	rootFromDot := yp.Get(".")
	dotMap, dotOk := rootFromDot.(map[string]interface{})

	if !dotOk {
		t.Error("Get(\".\") should return map[string]interface{}")
	} else if len(dataMap) != len(dotMap) {
		t.Error("Get() and Get(\".\") should return maps of same size")
	}

	// Test that returned data contains expected keys
	if _, ok := dataMap["server"]; !ok {
		t.Error("Get() should contain 'server' key")
	}

	if _, ok := dataMap["database"]; !ok {
		t.Error("Get() should contain 'database' key")
	}
}

func TestExists(t *testing.T) {
	yp := createSampleYPath(t)

	tests := []struct {
		path     string
		expected bool
	}{
		{"server.host", true},
		{"server.port", true},
		{"database.connection.pool.max", true},
		{"nonexistent.path", false},
		{"", true},
		{".", true},
	}

	for _, test := range tests {
		result := yp.Exists(test.path)
		if result != test.expected {
			t.Errorf("Exists(%s) = %t, expected %t", test.path, result, test.expected)
		}
	}
}
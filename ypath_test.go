package ypath

import (
	"os"
	"strings"
	"testing"
)

const sampleYAML = `
server:
  host: localhost
  port: 8080
  ssl: true
database:
  host: db.example.com
  port: 5432
  connection:
    pool:
      max: 10
      min: 2
    timeout: 30.5
  credentials:
    username: admin
    password: secret123
features:
  - authentication
  - logging
  - monitoring
servers:
  - host: server1.example.com
    port: 8001
    active: true
  - host: server2.example.com
    port: 8002
    active: false
numbers:
  int_val: 42
  float_val: 3.14
  string_num: "123"
  bool_true: "true"
  bool_false: "false"
  bool_one: "1"
`

func TestNew(t *testing.T) {
	yp, err := NewFromString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}
	if yp == nil {
		t.Fatal("YPath instance is nil")
	}
}

func TestGetString(t *testing.T) {
	yp, err := NewFromString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}

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
	yp, err := NewFromString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}

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
	yp, err := NewFromString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}

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
	yp, err := NewFromString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}

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

func TestArrayAccess(t *testing.T) {
	yp, err := NewFromString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}

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

func TestGet(t *testing.T) {
	yp, err := NewFromString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}

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

func TestExists(t *testing.T) {
	yp, err := NewFromString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}

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

func TestGetAllWildcard(t *testing.T) {
	yp, err := NewFromString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}

	// Test getting all features
	features := yp.GetAll("features.*")
	if len(features) != 3 {
		t.Errorf("GetAll(features.*) returned %d items, expected 3", len(features))
	}

	expectedFeatures := []string{"authentication", "logging", "monitoring"}
	for i, feature := range features {
		if feature != expectedFeatures[i] {
			t.Errorf("GetAll(features.*)[%d] = %v, expected %v", i, feature, expectedFeatures[i])
		}
	}

	// Test getting all server hosts
	serverHosts := yp.GetAll("servers.*.host")
	if len(serverHosts) != 2 {
		t.Errorf("GetAll(servers.*.host) returned %d items, expected 2", len(serverHosts))
	}

	expectedHosts := []string{"server1.example.com", "server2.example.com"}
	for i, host := range serverHosts {
		if host != expectedHosts[i] {
			t.Errorf("GetAll(servers.*.host)[%d] = %v, expected %v", i, host, expectedHosts[i])
		}
	}

	// Test non-wildcard path (should work like Get)
	singleValue := yp.GetAll("server.host")
	if len(singleValue) != 1 || singleValue[0] != "localhost" {
		t.Errorf("GetAll(server.host) = %v, expected [localhost]", singleValue)
	}
}

func TestEmptyYAML(t *testing.T) {
	yp, err := NewFromString("")
	if err != nil {
		t.Fatalf("Failed to create YPath from empty string: %v", err)
	}

	// Test operations on empty YAML
	if val := yp.GetString("any.path"); val != "" {
		t.Errorf("GetString on empty YAML returned %q, expected empty string", val)
	}

	if val := yp.GetInt("any.path"); val != 0 {
		t.Errorf("GetInt on empty YAML returned %d, expected 0", val)
	}

	if val := yp.GetBool("any.path"); val != false {
		t.Errorf("GetBool on empty YAML returned %t, expected false", val)
	}

	if val := yp.GetFloat64("any.path"); val != 0.0 {
		t.Errorf("GetFloat64 on empty YAML returned %f, expected 0.0", val)
	}
}

func TestInvalidYAML(t *testing.T) {
	invalidYAML := `
server:
  host: localhost
  port: 8080
invalid: [unclosed array
`

	_, err := NewFromString(invalidYAML)
	if err == nil {
		t.Error("Expected error when parsing invalid YAML, but got nil")
	}
}

func TestNewFromFile(t *testing.T) {
	// Create a temporary YAML file for testing
	tempFile := "test_config.yaml"
	testYAML := `
server:
  host: test.example.com
  port: 9999
  ssl: false
database:
  host: db.test.com
  connection:
    pool:
      max: 20
`

	// Write test YAML to file
	err := os.WriteFile(tempFile, []byte(testYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() {
		// Clean up test file
		if removeErr := os.Remove(tempFile); removeErr != nil {
			t.Logf("Warning: failed to remove test file: %v", removeErr)
		}
	}()

	// Test NewFromFile
	yp, err := NewFromFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to create YPath from file: %v", err)
	}

	// Verify values loaded correctly
	if host := yp.GetString("server.host"); host != "test.example.com" {
		t.Errorf("Expected server.host to be 'test.example.com', got '%s'", host)
	}

	if port := yp.GetInt("server.port"); port != 9999 {
		t.Errorf("Expected server.port to be 9999, got %d", port)
	}

	if ssl := yp.GetBool("server.ssl"); ssl != false {
		t.Errorf("Expected server.ssl to be false, got %t", ssl)
	}

	if maxPool := yp.GetInt("database.connection.pool.max"); maxPool != 20 {
		t.Errorf("Expected database.connection.pool.max to be 20, got %d", maxPool)
	}
}

func TestNewFromFile_NonExistent(t *testing.T) {
	_, err := NewFromFile("nonexistent_file.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent file, but got nil")
	}

	// Verify error message contains file name
	expectedErr := "failed to read file nonexistent_file.yaml"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestNewFromFile_InvalidYAML(t *testing.T) {
	// Create a temporary file with invalid YAML
	tempFile := "invalid_test.yaml"
	invalidYAML := `
server:
  host: localhost
  invalid: [unclosed array
`

	err := os.WriteFile(tempFile, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid test file: %v", err)
	}
	defer func() {
		if removeErr := os.Remove(tempFile); removeErr != nil {
			t.Logf("Warning: failed to remove test file: %v", removeErr)
		}
	}()

	_, err = NewFromFile(tempFile)
	if err == nil {
		t.Error("Expected error when parsing invalid YAML file, but got nil")
	}

	// Verify error message indicates YAML parsing error
	if !strings.Contains(err.Error(), "failed to parse YAML") {
		t.Errorf("Expected error message to contain 'failed to parse YAML', got '%s'", err.Error())
	}
}
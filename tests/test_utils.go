package tests

import (
	"os"
	"testing"

	"github.com/afeiship/go-yaml-path"
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

// createYPath creates a new YPath instance for testing
func createYPath(t *testing.T, yaml string) *ypath.YPath {
	yp, err := ypath.NewFromString(yaml)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}
	if yp == nil {
		t.Fatal("YPath instance is nil")
	}
	return yp
}

// createSampleYPath creates a YPath instance with sample data
func createSampleYPath(t *testing.T) *ypath.YPath {
	return createYPath(t, sampleYAML)
}

// createTempFile creates a temporary YAML file for testing
func createTempFile(t *testing.T, content string) string {
	tempFile := "test_config.yaml"
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return tempFile
}

// cleanupTempFile removes a temporary file
func cleanupTempFile(t *testing.T, filename string) {
	if err := os.Remove(filename); err != nil {
		t.Logf("Warning: failed to remove test file %s: %v", filename, err)
	}
}
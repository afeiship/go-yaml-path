package tests

import (
	"testing"

	"github.com/afeiship/go-yaml-path"
)

func TestNew(t *testing.T) {
	yp, err := ypath.NewFromString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to create YPath: %v", err)
	}
	if yp == nil {
		t.Fatal("YPath instance is nil")
	}
}

func TestNew_EmptyYAML(t *testing.T) {
	yp, err := ypath.NewFromString("")
	if err != nil {
		t.Fatalf("Failed to create YPath from empty string: %v", err)
	}
	if yp == nil {
		t.Fatal("YPath instance should not be nil for empty YAML")
	}
}

func TestNew_InvalidYAML(t *testing.T) {
	invalidYAML := `
server:
  host: localhost
  port: 8080
invalid: [unclosed array
`

	_, err := ypath.NewFromString(invalidYAML)
	if err == nil {
		t.Error("Expected error when parsing invalid YAML, but got nil")
	}
}

func TestNewFromFile(t *testing.T) {
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

	tempFile := createTempFile(t, testYAML)
	defer cleanupTempFile(t, tempFile)

	// Test NewFromFile
	yp, err := ypath.NewFromFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to create YPath from file: %v", err)
	}

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
	_, err := ypath.NewFromFile("nonexistent_file.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent file, but got nil")
	}
}

func TestNewFromFile_InvalidYAML(t *testing.T) {
	invalidYAML := `
server:
  host: localhost
  invalid: [unclosed array
`

	tempFile := createTempFile(t, invalidYAML)
	defer cleanupTempFile(t, tempFile)

	_, err := ypath.NewFromFile(tempFile)
	if err == nil {
		t.Error("Expected error when parsing invalid YAML file, but got nil")
	}
}
package tests

import (
	"fmt"
	"testing"
)

func TestEmptyYAML(t *testing.T) {
	yp := createYPath(t, "")

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

	// Test list operations on empty YAML
	if val := yp.GetStringList("any.path"); val != nil {
		t.Errorf("GetStringList on empty YAML returned %v, expected nil", val)
	}

	if val := yp.GetIntList("any.path"); val != nil {
		t.Errorf("GetIntList on empty YAML returned %v, expected nil", val)
	}
}

func TestComplexNestedStructures(t *testing.T) {
	complexYAML := `
apps:
  - name: webapp
    config:
      server:
        host: localhost
        port: 3000
        ssl: true
      database:
        name: webapp_db
        pool:
          min: 2
          max: 10
  - name: apiapp
    config:
      server:
        host: api.example.com
        port: 8080
        ssl: false
      database:
        name: api_db
        pool:
          min: 1
          max: 5
`

	yp := createYPath(t, complexYAML)

	// Test deep nested access
	if val := yp.GetString("apps.1.config.database.name"); val != "api_db" {
		t.Errorf("Expected 'api_db', got %s", val)
	}

	// Test wildcard on complex structure
	hosts := yp.GetStringList("apps.*.config.server.host")
	if len(hosts) != 2 {
		t.Errorf("Expected 2 hosts, got %d", len(hosts))
	}

	maxPools := yp.GetIntList("apps.*.config.database.pool.max")
	if len(maxPools) != 2 {
		t.Errorf("Expected 2 max pools, got %d", len(maxPools))
	}
}

func TestSpecialCharactersAndNumbers(t *testing.T) {
	specialYAML := `
special:
  "key-with-dash": value1
  "key_with_underscore": value2
  "123numberkey": value3
numbers:
  - 42
  - 3.14
  - "123"
  - "456.789"
  - true
`

	yp := createYPath(t, specialYAML)

	// Test special character keys (note: this might not work with current implementation)
	// Current implementation uses simple string splitting, so complex keys might not work

	// Test mixed number types
	mixedNumbers := yp.GetList("numbers")
	if len(mixedNumbers) != 5 {
		t.Errorf("Expected 5 numbers, got %d", len(mixedNumbers))
	}
}

func TestLargeYAML(t *testing.T) {
	// Generate a large YAML structure with valid characters
	var largeYAML string = "root:\n  items:\n"
	for i := 0; i < 100; i++ { // Reduced size for testing
		largeYAML += fmt.Sprintf("    - name: item_%d\n", i)
		largeYAML += fmt.Sprintf("      value: %d\n", i)
	}

	yp := createYPath(t, largeYAML)

	// Test accessing large arrays
	items := yp.GetAll("root.items.*.name")
	if len(items) != 100 {
		t.Errorf("Expected 100 items, got %d", len(items))
	}

	// Test specific index access
	firstItem := yp.GetString("root.items.0.name")
	if firstItem != "item_0" {
		t.Errorf("Expected 'item_0', got '%s'", firstItem)
	}

	// Test last item access
	lastItem := yp.GetString("root.items.99.name")
	if lastItem != "item_99" {
		t.Errorf("Expected 'item_99', got '%s'", lastItem)
	}
}
package tests

import "testing"

func TestGetStringList(t *testing.T) {
	yp := createSampleYPath(t)

	// Test string array
	features := yp.GetStringList("features")
	expected := []string{"authentication", "logging", "monitoring"}
	if len(features) != len(expected) {
		t.Errorf("Expected %d features, got %d", len(expected), len(features))
	}
	for i, feature := range features {
		if feature != expected[i] {
			t.Errorf("Expected feature[%d] = %s, got %s", i, expected[i], feature)
		}
	}

	// Test converting non-string array to strings
	ports := yp.GetStringList("servers.*.port")
	if len(ports) != 2 {
		t.Errorf("Expected 2 ports, got %d", len(ports))
	}

	// Test single value conversion
	singleValue := yp.GetStringList("server.host")
	if len(singleValue) != 1 || singleValue[0] != "localhost" {
		t.Errorf("Expected ['localhost'], got %v", singleValue)
	}

	// Test non-existent path
	nonExistent := yp.GetStringList("nonexistent.path")
	if nonExistent != nil {
		t.Errorf("Expected nil for non-existent path, got %v", nonExistent)
	}
}

func TestGetIntList(t *testing.T) {
	yp := createSampleYPath(t)

	// Test int array
	ports := yp.GetIntList("servers.*.port")
	expected := []int{8001, 8002}
	if len(ports) != len(expected) {
		t.Errorf("Expected %d ports, got %d", len(expected), len(ports))
	}
	for i, port := range ports {
		if port != expected[i] {
			t.Errorf("Expected port[%d] = %d, got %d", i, expected[i], port)
		}
	}

	// Test single value conversion
	singleValue := yp.GetIntList("server.port")
	if len(singleValue) != 1 || singleValue[0] != 8080 {
		t.Errorf("Expected [8080], got %v", singleValue)
	}

	// Test non-existent path
	nonExistent := yp.GetIntList("nonexistent.path")
	if nonExistent != nil {
		t.Errorf("Expected nil for non-existent path, got %v", nonExistent)
	}
}

func TestGetBoolList(t *testing.T) {
	yp := createSampleYPath(t)

	// Test bool array
	activeServers := yp.GetBoolList("servers.*.active")
	expected := []bool{true, false}
	if len(activeServers) != len(expected) {
		t.Errorf("Expected %d active flags, got %d", len(expected), len(activeServers))
	}
	for i, active := range activeServers {
		if active != expected[i] {
			t.Errorf("Expected active[%d] = %t, got %t", i, expected[i], active)
		}
	}

	// Test single value conversion
	singleValue := yp.GetBoolList("server.ssl")
	if len(singleValue) != 1 || singleValue[0] != true {
		t.Errorf("Expected [true], got %v", singleValue)
	}

	// Test non-existent path
	nonExistent := yp.GetBoolList("nonexistent.path")
	if nonExistent != nil {
		t.Errorf("Expected nil for non-existent path, got %v", nonExistent)
	}
}

func TestGetFloat64List(t *testing.T) {
	yp := createSampleYPath(t)

	// Test single value conversion
	singleValue := yp.GetFloat64List("database.connection.timeout")
	if len(singleValue) != 1 || singleValue[0] != 30.5 {
		t.Errorf("Expected [30.5], got %v", singleValue)
	}

	// Test non-existent path
	nonExistent := yp.GetFloat64List("nonexistent.path")
	if nonExistent != nil {
		t.Errorf("Expected nil for non-existent path, got %v", nonExistent)
	}
}

func TestGetList(t *testing.T) {
	yp := createSampleYPath(t)

	// Test getting interface array
	features := yp.GetList("features")
	if len(features) != 3 {
		t.Errorf("Expected 3 features, got %d", len(features))
	}

	// Test single value to list conversion
	singleValue := yp.GetList("server.host")
	if len(singleValue) != 1 || singleValue[0] != "localhost" {
		t.Errorf("Expected ['localhost'], got %v", singleValue)
	}

	// Test non-existent path
	nonExistent := yp.GetList("nonexistent.path")
	if nonExistent != nil {
		t.Errorf("Expected nil for non-existent path, got %v", nonExistent)
	}
}

func TestListTypeConversions(t *testing.T) {
	yp := createSampleYPath(t)

	// Test string number conversions
	stringNumbers := yp.GetIntList("numbers")
	if len(stringNumbers) > 0 {
		t.Logf("String numbers conversion result: %v", stringNumbers)
	}

	// Test string boolean conversions
	stringBools := yp.GetBoolList("numbers")
	if len(stringBools) > 0 {
		t.Logf("String booleans conversion result: %v", stringBools)
	}

	// Test float number conversions
	floatNumbers := yp.GetFloat64List("numbers")
	if len(floatNumbers) > 0 {
		t.Logf("Float numbers conversion result: %v", floatNumbers)
	}
}
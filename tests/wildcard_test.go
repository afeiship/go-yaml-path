package tests

import "testing"

func TestGetAllWildcard(t *testing.T) {
	yp := createSampleYPath(t)

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

func TestWildcardMultipleLevels(t *testing.T) {
	yp := createSampleYPath(t)

	// Test wildcard at different levels
	allPorts := yp.GetAll("servers.*.port")
	if len(allPorts) != 2 {
		t.Errorf("Expected 2 ports, got %d", len(allPorts))
	}

	// Test nested wildcard
	poolValues := yp.GetAll("database.connection.pool.*")
	if len(poolValues) < 2 {
		t.Errorf("Expected at least 2 pool values, got %d", len(poolValues))
	}
}

func TestWildcardNonExistentPath(t *testing.T) {
	yp := createSampleYPath(t)

	// Test wildcard on non-existent path
	result := yp.GetAll("nonexistent.*.path")
	if result != nil {
		t.Errorf("Expected nil for non-existent wildcard path, got %v", result)
	}
}

func TestWildcardEmptyArrays(t *testing.T) {
	emptyYAML := `
empty_array: []
servers: []
features: []
`

	yp := createYPath(t, emptyYAML)

	// Test wildcard on empty arrays
	result := yp.GetAll("empty_array.*")
	if len(result) != 0 {
		t.Errorf("Expected empty array for wildcard on empty array, got %v", result)
	}

	result2 := yp.GetAll("servers.*.host")
	if len(result2) != 0 {
		t.Errorf("Expected empty array for wildcard on empty servers, got %v", result2)
	}
}
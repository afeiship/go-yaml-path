package main

import (
	"fmt"
	"log"

	"github.com/afeiship/go-yaml-path"
)

func main() {
	// Sample YAML configuration
	yamlConfig := `
server:
  host: localhost
  port: 8080
  ssl: true
  api:
    version: v1
    endpoint: /api/v1

database:
  host: db.example.com
  port: 5432
  name: myapp
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
  - caching

servers:
  - host: server1.example.com
    port: 8001
    active: true
    region: us-east-1
  - host: server2.example.com
    port: 8002
    active: false
    region: us-west-1
  - host: server3.example.com
    port: 8003
    active: true
    region: eu-west-1

metrics:
  requests_per_second: 1000
  error_rate: 0.01
  uptime_percentage: 99.9
`

	// Create YPath instance from string
	yp, err := ypath.NewFromString(yamlConfig)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	fmt.Println("=== YAML Path Library Demo ===")
	fmt.Println()

	// Basic string values
	fmt.Println("Basic String Values:")
	fmt.Printf("Server Host: %s\n", yp.GetString("server.host"))
	fmt.Printf("Database Host: %s\n", yp.GetString("database.host"))
	fmt.Printf("API Version: %s\n", yp.GetString("server.api.version"))
	fmt.Printf("Database Name: %s\n", yp.GetString("database.name"))
	fmt.Println()

	// Integer values
	fmt.Println("Integer Values:")
	fmt.Printf("Server Port: %d\n", yp.GetInt("server.port"))
	fmt.Printf("Database Port: %d\n", yp.GetInt("database.port"))
	fmt.Printf("Max Pool Size: %d\n", yp.GetInt("database.connection.pool.max"))
	fmt.Printf("Min Pool Size: %d\n", yp.GetInt("database.connection.pool.min"))
	fmt.Println()

	// Float values
	fmt.Println("Float Values:")
	fmt.Printf("Connection Timeout: %.1f\n", yp.GetFloat64("database.connection.timeout"))
	fmt.Printf("Error Rate: %.2f%%\n", yp.GetFloat64("metrics.error_rate")*100)
	fmt.Printf("Uptime: %.1f%%\n", yp.GetFloat64("metrics.uptime_percentage"))
	fmt.Println()

	// Boolean values
	fmt.Println("Boolean Values:")
	fmt.Printf("SSL Enabled: %t\n", yp.GetBool("server.ssl"))
	fmt.Printf("Server 1 Active: %t\n", yp.GetBool("servers.0.active"))
	fmt.Printf("Server 2 Active: %t\n", yp.GetBool("servers.1.active"))
	fmt.Printf("Server 3 Active: %t\n", yp.GetBool("servers.2.active"))
	fmt.Println()

	// Array access
	fmt.Println("Array Access:")
	fmt.Printf("First Feature: %s\n", yp.GetString("features.0"))
	fmt.Printf("Second Feature: %s\n", yp.GetString("features.1"))
	fmt.Printf("Third Feature: %s\n", yp.GetString("features.2"))
	fmt.Printf("Fourth Feature: %s\n", yp.GetString("features.3"))
	fmt.Println()

	// Nested array access
	fmt.Println("Nested Array Access:")
	fmt.Printf("First Server Host: %s\n", yp.GetString("servers.0.host"))
	fmt.Printf("First Server Port: %d\n", yp.GetInt("servers.0.port"))
	fmt.Printf("First Server Region: %s\n", yp.GetString("servers.0.region"))
	fmt.Printf("Second Server Host: %s\n", yp.GetString("servers.1.host"))
	fmt.Printf("Third Server Region: %s\n", yp.GetString("servers.2.region"))
	fmt.Println()

	// Check if paths exist
	fmt.Println("Path Existence Checks:")
	fmt.Printf("Path 'server.host' exists: %t\n", yp.Exists("server.host"))
	fmt.Printf("Path 'server.ssl' exists: %t\n", yp.Exists("server.ssl"))
	fmt.Printf("Path 'nonexistent.path' exists: %t\n", yp.Exists("nonexistent.path"))
	fmt.Println()

	// Wildcard operations
	fmt.Println("Wildcard Operations:")
	allFeatures := yp.GetAll("features.*")
	fmt.Printf("All Features: %v\n", allFeatures)

	allServerHosts := yp.GetAll("servers.*.host")
	fmt.Printf("All Server Hosts: %v\n", allServerHosts)

	allActiveServers := yp.GetAll("servers.*.active")
	fmt.Printf("All Server Active Status: %v\n", allActiveServers)

	allRegions := yp.GetAll("servers.*.region")
	fmt.Printf("All Server Regions: %v\n", allRegions)
	fmt.Println()

	// Get raw values
	fmt.Println("Raw Value Access:")
	serverConfig := yp.Get("server")
	fmt.Printf("Server Config: %v\n", serverConfig)

	databaseConfig := yp.Get("database.connection")
	fmt.Printf("Database Connection Config: %v\n", databaseConfig)

	features := yp.Get("features")
	fmt.Printf("Features Array: %v\n", features)
	fmt.Println()

	// Demonstrate Get() method without parameters
	fmt.Println("=== Get() Method Demo (No Parameters) ===")
	rootData := yp.Get()
	fmt.Printf("Root Data (using Get()): %T\n", rootData)

	// Compare Get() without parameters with Get("")
	dataFromEmpty := yp.Get("")
	dataMap, emptyOk := dataFromEmpty.(map[string]any)
	rootMap, rootOk := rootData.(map[string]any)

	if emptyOk && rootOk && len(dataMap) == len(rootMap) {
		fmt.Println("Get() == Get(\"\"): true (same size and type)")
	} else {
		fmt.Println("Get() == Get(\"\"): false")
	}

	// Handle non-existent paths gracefully
	fmt.Println("Handling Non-existent Paths:")
	fmt.Printf("Non-existent string: '%s'\n", yp.GetString("nonexistent.string"))
	fmt.Printf("Non-existent int: %d\n", yp.GetInt("nonexistent.int"))
	fmt.Printf("Non-existent bool: %t\n", yp.GetBool("nonexistent.bool"))
	fmt.Printf("Non-existent float: %f\n", yp.GetFloat64("nonexistent.float"))
	fmt.Printf("Non-existent raw value: %v\n", yp.Get("nonexistent.path"))
	fmt.Println()

	fmt.Println()

	// Demonstrate loading from file
	fmt.Println("=== Loading from File ===")
	ypFile, err := ypath.NewFromFile("config.yaml")
	if err != nil {
		log.Printf("Failed to load config.yaml: %v", err)
	} else {
		fmt.Printf("Loaded config from file:\n")
		fmt.Printf("Server: %s:%d\n", ypFile.GetString("server.host"), ypFile.GetInt("server.port"))
		fmt.Printf("Database: %s:%d\n", ypFile.GetString("database.host"), ypFile.GetInt("database.port"))
		fmt.Printf("SSL Enabled: %t\n", ypFile.GetBool("server.ssl"))
		fmt.Printf("Features: %v\n", ypFile.GetAll("features.*"))
	}

	fmt.Println()

	// Demonstrate GetList methods
	fmt.Println("=== GetList Methods Demo ===")

	// String list
	featuresList := yp.GetStringList("features")
	fmt.Printf("Features (string list): %v\n", featuresList)

	// Integer list with wildcard
	ports := yp.GetIntList("servers.*.port")
	fmt.Printf("Server Ports (int list with wildcard): %v\n", ports)

	// Boolean list with wildcard
	activeServers := yp.GetBoolList("servers.*.active")
	fmt.Printf("Active Servers (bool list): %v\n", activeServers)

	// Generic list
	serverList := yp.GetList("servers")
	fmt.Printf("Server Configs (generic list): %v\n", serverList)

	// Single value to list conversion
	singleHost := yp.GetStringList("server.host")
	fmt.Printf("Single Host as List: %v\n", singleHost)

	fmt.Println("\n=== Demo Complete ===")
}
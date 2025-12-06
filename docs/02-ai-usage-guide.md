# Go-YAML-Path AI Usage Guide

This guide provides comprehensive documentation for AI systems to understand and use the go-yaml-path library effectively.

## Overview

The go-yaml-path library is a Go package that provides intuitive dot notation path navigation for YAML data without requiring struct definitions. It enables developers to access nested YAML values using familiar JavaScript-like dot syntax with additional features like wildcard matching and type-safe accessors.

## Quick Start

### Installation
```bash
go get github.com/feizheng/go-yaml-path
```

### Basic Usage
```go
package main

import (
    "fmt"
    "github.com/feizheng/go-yaml-path/ypath"
)

func main() {
    yamlData := `
server:
  host: localhost
  port: 8080
  ssl: true
database:
  connection:
    pool:
      max: 10
features:
  - authentication
  - logging
  - monitoring
servers:
  - host: server1.example.com
    port: 8001
    region: us-east-1
  - host: server2.example.com
    port: 8002
    region: us-west-2
`

    yp, err := ypath.NewFromString(yamlData)
    if err != nil {
        panic(err)
    }

    // Access configuration values
    host := yp.GetString("server.host")      // "localhost"
    port := yp.GetInt("server.port")         // 8080
    ssl := yp.GetBool("server.ssl")         // true
    maxPool := yp.GetInt("database.connection.pool.max") // 10

    fmt.Printf("Host: %s, Port: %d, SSL: %t, Max Pool: %d\n",
        host, port, ssl, maxPool)
}
```

## Core API

### Factory Functions

#### `New(yamlData []byte) (*YPath, error)`
Creates a new YPath instance from YAML byte data.

#### `NewFromString(yamlStr string) (*YPath, error)`
Creates a new YPath instance from YAML string.

#### `NewFromFile(filename string) (*YPath, error)`
Creates a new YPath instance from a YAML file.

### Core Accessor Methods

#### `Get(path ...string) interface{}`
Returns the raw value at the specified path.
- If no path is provided, returns the entire data structure
- Returns `nil` for non-existent paths
- Supports nested access with dot notation

#### Type-Specific Getters
- `GetString(path string) string` - Returns string value with type conversion
- `GetInt(path string) int` - Returns integer value with type conversion
- `GetBool(path string) bool` - Returns boolean value with type conversion
- `GetFloat64(path string) float64` - Returns float64 value with type conversion
- `Exists(path string) bool` - Returns true if path exists

### Wildcard and List Methods

#### `GetAll(path string) []interface{}`
Returns all values matching a wildcard path.

#### Type-Specific List Methods
- `GetStringList(path string) []string` - String list with auto-conversion
- `GetIntList(path string) []int` - Integer list with type conversion
- `GetBoolList(path string) []bool` - Boolean list with conversion
- `GetFloat64List(path string) []float64` - Float64 list with conversion
- `GetList(path string) []interface{}` - Generic list method

## Path Expression Syntax

### 1. Basic Dot Notation
```
server.host → Access "host" key in "server" object
database.connection.pool.max → Deeply nested access
api.version → Multi-level nested access
```

### 2. Array Indexing
```
features.0 → First element of "features" array
servers.1.host → "host" property of second server
items.2.value → Value at third item
```

### 3. Wildcard Support
```
servers.*.host → All server hosts
features.* → All feature names
database.*.max → All "max" values under database
config.*.enabled → All "enabled" flags
```

### 4. Special Cases
- `""` (empty string) or `"."` → Returns entire YAML data structure
- Non-existent paths → Return appropriate zero values (no errors)

## Type Conversion Rules

### GetString()
- Any type → string using `fmt.Sprintf("%v")`
- Numbers: `123` → `"123"`
- Booleans: `true` → `"true"`
- Objects/Arrays: Complex string representation

### GetInt()
- `int`, `int64` → Direct conversion
- `float64` → Truncated to int (123.9 → 123)
- String → Parsed as integer if possible
- Invalid conversion → Returns 0

### GetBool()
- `bool` → Direct value
- String → `"true"` or `"1"` → `true`, others → `false`
- Numbers → Non-zero → `true`, zero → `false`
- Invalid conversion → Returns `false`

### GetFloat64()
- `float64` → Direct value
- `int`, `int64` → Converted to float64
- String → Parsed as float64 if possible
- Invalid conversion → Returns 0.0

## Usage Patterns

### Configuration Management
```go
yp, err := ypath.NewFromFile("config.yaml")
if err != nil {
    log.Fatal(err)
}

// Database configuration
dbHost := yp.GetString("database.host")
dbPort := yp.GetInt("database.port")
dbSSL := yp.GetBool("database.ssl")

// API settings
apiTimeout := yp.GetFloat64("api.timeout")
apiRetries := yp.GetInt("api.retries")
```

### Kubernetes YAML Processing
```go
yp, err := ypath.NewFromFile("deployment.yaml")

// Extract pod information
podName := yp.GetString("metadata.name")
replicas := yp.GetInt("spec.replicas")

// Get container names
containerNames := yp.GetStringList("spec.template.spec.containers.*.name")

// Get all container images
containerImages := yp.GetStringList("spec.template.spec.containers.*.image")
```

### Batch Operations with Wildcards
```go
// Get all server configurations
servers := yp.GetAll("servers.*")
for i, server := range servers {
    fmt.Printf("Server %d: %v\n", i, server)
}

// Extract specific properties
hosts := yp.GetStringList("servers.*.host")
ports := yp.GetIntList("servers.*.port")

// Check all service statuses
statuses := yp.GetBoolList("services.*.enabled")
```

### Error Handling Best Practices
```go
// Check existence before accessing critical paths
if yp.Exists("database.host") {
    dbHost := yp.GetString("database.host")
    // Process dbHost
}

// Use zero-value checks for optional configuration
port := yp.GetInt("server.port")
if port == 0 {
    port = 8080 // Use default
}

// Type-safe processing
apiTimeout := yp.GetFloat64("api.timeout")
if apiTimeout == 0.0 {
    log.Warn("API timeout not configured, using default")
    apiTimeout = 30.0
}
```

## Advanced Examples

### Complex Nested Structure Access
```go
yamlData := `
environments:
  production:
    servers:
      - host: prod1.example.com
        config:
          resources:
            cpu: 4
            memory: "8Gi"
      - host: prod2.example.com
        config:
          resources:
            cpu: 2
            memory: "4Gi"
  staging:
    servers:
      - host: staging.example.com
        config:
          resources:
            cpu: 1
            memory: "2Gi"
`

yp, _ := ypath.NewFromString(yamlData)

// Access deeply nested values
prodCPU := yp.GetInt("environments.production.servers.0.config.resources.cpu")
prodMem := yp.GetString("environments.production.servers.0.config.resources.memory")

// Get all CPU configurations across environments
allCPUs := yp.GetIntList("environments.*.servers.*.config.resources.cpu")

// Get all production server hosts
prodHosts := yp.GetStringList("environments.production.servers.*.host")
```

### Conditional Access and Defaults
```go
// Helper function to get values with defaults
func getStringWithDefault(yp *ypath.YPath, path, defaultValue string) string {
    if yp.Exists(path) {
        value := yp.GetString(path)
        if value != "" {
            return value
        }
    }
    return defaultValue
}

// Usage
host := getStringWithDefault(yp, "server.host", "localhost")
port := getIntWithDefault(yp, "server.port", 8080)
```

## Performance Considerations

- **Time Complexity**: O(n) where n is the path depth
- **Wildcard Operations**: O(m×n) where m is the number of matching elements
- **Memory Usage**: Minimal - only stores parsed YAML data
- **Type Conversion**: Adds minimal overhead
- **No Caching**: Values are recomputed on each access

## Common Use Cases

### 1. Application Configuration
```go
yp, _ := ypath.NewFromFile("app.yaml")
appPort := yp.GetInt("server.port")
dbURL := yp.GetString("database.url")
logLevel := yp.GetString("logging.level")
```

### 2. CI/CD Pipeline Configuration
```go
yp, _ := ypath.NewFromFile("pipeline.yml")
stages := yp.GetStringList("pipeline.stages.*.name")
timeout := yp.GetInt("pipeline.global.timeout")
```

### 3. Test Data Management
```go
yp, _ := ypath.NewFromFile("test-data.yaml")
testUsers := yp.GetStringList("users.*.username")
testPasswords := yp.GetStringList("users.*.password")
```

### 4. API Response Processing
```go
yp, _ := ypath.New(respData)
userID := yp.GetString("data.user.id")
permissions := yp.GetStringList("data.user.permissions.*")
```

## Integration Tips for AI Systems

### 1. Code Generation
When generating code with go-yaml-path:
- Use type-specific getters for type safety
- Include existence checks for critical paths
- Handle empty returns gracefully

### 2. Configuration Validation
```go
func validateConfig(yp *ypath.YPath) []string {
    var errors []string

    requiredPaths := []string{
        "server.host",
        "server.port",
        "database.url",
    }

    for _, path := range requiredPaths {
        if !yp.Exists(path) {
            errors = append(errors, fmt.Sprintf("Missing required path: %s", path))
        }
    }

    return errors
}
```

### 3. Dynamic Path Construction
```go
// Build paths dynamically for flexible configuration
func getValueByPrefix(yp *ypath.YPath, prefix, key string) string {
    path := fmt.Sprintf("%s.%s", prefix, key)
    return yp.GetString(path)
}

// Usage
dbHost := getValueByPrefix(yp, "database", "host")
dbPort := getIntWithDefault(yp, "database.port", 5432)
```

## Error Handling

The library follows a "fail-safe" approach:
- Non-existent paths return zero values rather than errors
- Invalid type conversions return zero values
- No panics are raised from normal operations
- Errors only occur during YAML parsing

## Frequently Used Patterns

### Reading Environment-Specific Configs
```go
func getEnvConfig(yp *ypath.YPath, env string) map[string]interface{} {
    config := make(map[string]interface{})

    // Get environment-specific overrides
    envPath := fmt.Sprintf("environments.%s", env)
    if yp.Exists(envPath) {
        envData := yp.Get(envPath)
        if envMap, ok := envData.(map[string]interface{}); ok {
            return envMap
        }
    }

    return config
}
```

### Processing Arrays of Objects
```go
// Extract specific fields from object arrays
func extractFields(yp *ypath.YPath, arrayPath, field string) []string {
    values := yp.GetStringList(fmt.Sprintf("%s.*.%s", arrayPath, field))
    return values
}

// Usage
serverHosts := extractFields(yp, "servers", "host")
containerPorts := extractFields(yp, "containers", "port")
```

This guide provides AI systems with comprehensive knowledge to effectively utilize the go-yaml-path library in various scenarios.
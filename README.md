# go-yaml-path

A Go library for accessing YAML data using dot notation paths without defining structs.

**[中文文档](README.zh-CN.md) | English**

## Features

- ✅ Uses `gopkg.in/yaml.v3` for YAML parsing
- ✅ No struct definitions required
- ✅ Provides `GetString()`, `GetInt()`, `GetBool()`, `GetFloat64()` methods
- ✅ Supports dot syntax for nested value access
- ✅ Array index access (e.g., `servers.0.host`)
- ✅ Wildcard support (e.g., `servers.*.host`)
- ✅ Type conversion and safe defaults
- ✅ Path existence checking

## Installation

```bash
go get github.com/afeiship/go-yaml-path
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/afeiship/go-yaml-path"
)

func main() {
    yamlData := `
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
features:
  - authentication
  - logging
  - monitoring
`

    // Create YPath instance
    yp, err := ypath.NewFromString(yamlData)
    if err != nil {
        panic(err)
    }

    // Access values using dot notation
    fmt.Println("Server Host:", yp.GetString("server.host"))        // localhost
    fmt.Println("Server Port:", yp.GetInt("server.port"))           // 8080
    fmt.Println("SSL Enabled:", yp.GetBool("server.ssl"))           // true
    fmt.Println("Max Pool Size:", yp.GetInt("database.connection.pool.max")) // 10

    // Array access
    fmt.Println("First Feature:", yp.GetString("features.0"))       // authentication

    // Check if path exists
    if yp.Exists("server.host") {
        fmt.Println("Server host is configured")
    }

    // Wildcard access
    features := yp.GetAll("features.*")
    fmt.Println("All features:", features)                          // [authentication logging monitoring]
}
```

## API Reference

### Creating YPath Instance

```go
// From byte slice
yp, err := ypath.New(yamlBytes)

// From string
yp, err := ypath.NewFromString(yamlString)

// From file
yp, err := ypath.NewFromFile("config.yaml")
```

#### `NewFromFile(filename string) (*YPath, error)`
Creates a new YPath instance from a YAML file.

```go
yp, err := ypath.NewFromFile("config.yaml")
if err != nil {
    panic(err)
}
```

### Basic Methods

#### `Get(path string) interface{}`
Retrieves raw value at the specified path.

```go
value := yp.Get("server.host")
```

#### `GetString(path string) string`
Retrieves a string value, converts other types to string.

```go
host := yp.GetString("server.host")  // "localhost"
port := yp.GetString("server.port")  // "8080" (converted from int)
```

#### `GetInt(path string) int`
Retrieves an integer value, with type conversion.

```go
port := yp.GetInt("server.port")           // 8080
maxPool := yp.GetInt("database.pool.max")  // 10
```

#### `GetBool(path string) bool`
Retrieves a boolean value, with string conversion.

```go
ssl := yp.GetBool("server.ssl")      // true
active := yp.GetBool("server.active") // false
```

#### `GetFloat64(path string) float64`
Retrieves a float64 value, with type conversion.

```go
timeout := yp.GetFloat64("connection.timeout")  // 30.5
```

#### `Exists(path string) bool`
Checks if a path exists in the YAML data.

```go
if yp.Exists("server.host") {
    // Path exists
}
```

### Advanced Methods

#### `GetAll(path string) []interface{}`
Retrieves all values matching a wildcard path.

```go
// Get all server hosts
hosts := yp.GetAll("servers.*.host")
// Result: []interface{}{"server1.example.com", "server2.example.com"}

// Get all features
features := yp.GetAll("features.*")
// Result: []interface{}{"authentication", "logging", "monitoring"}
```

## Path Syntax

### Dot Notation
Use dots to access nested values:
- `server.host` → access `host` in `server`
- `database.connection.pool.max` → deeply nested access

### Array Indexing
Use numeric indices to access array elements:
- `features.0` → first element of `features` array
- `servers.1.host` → `host` of second server

### Wildcards
Use `*` as wildcard for matching multiple values:
- `servers.*.host` → all server hosts
- `features.*` → all features
- `database.*.max` → all `max` values under `database`

### Special Cases
- Empty string `""` or dot `"."` → returns root YAML data
- Non-existent paths → return appropriate zero values (nil, "", 0, false)

## Type Conversion

The library automatically converts between types when possible:

- **GetString**: Any type → string using `fmt.Sprintf`
- **GetInt**: int, int64, float64 → int; strings → parsed as int if possible
- **GetBool**: bool → bool; strings → true for "true"/"1"; numbers → non-zero = true
- **GetFloat64**: float64, int, int64 → float64; strings → parsed as float64 if possible

## Error Handling

- **Invalid YAML**: Returns error from `yaml.Unmarshal`
- **Non-existent paths**: Return zero values, no error
- **Type conversion failures**: Return zero values, no error

## Example Use Cases

### Configuration Files
```go
yp, err := ypath.NewFromFile("config.yaml")
if err != nil {
    panic(err)
}
dbHost := yp.GetString("database.host")
dbPort := yp.GetInt("database.port")
enableSSL := yp.GetBool("server.ssl")
```

### Kubernetes Resources
```go
yp, _ := ypath.NewFromString(k8sYaml)
replicas := yp.GetInt("spec.replicas")
imageName := yp.GetString("spec.containers.0.image")
```

### API Responses
```go
yp, _ := ypath.NewFromString(apiResponse)
statusCode := yp.GetInt("response.status")
message := yp.GetString("response.message")
errors := yp.GetAll("response.errors.*.code")
```

## Running the Example

```bash
cd example
go run main.go
```

## Running Tests

```bash
go test -v
```

## License

MIT License
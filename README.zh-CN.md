# go-yaml-path

一个用于通过点号路径访问 YAML 数据的 Go 库，无需定义结构体。

**中文 | [English](README.md)**

## 特性

- ✅ 使用 `gopkg.in/yaml.v3` 解析 YAML
- ✅ 无需定义结构体
- ✅ 提供 `GetString()`、`GetInt()`、`GetBool()`、`GetFloat64()` 方法
- ✅ 支持点号语法访问嵌套值
- ✅ 数组索引访问（例如：`servers.0.host`）
- ✅ 通配符支持（例如：`servers.*.host`）
- ✅ 类型转换和安全默认值
- ✅ 路径存在性检查

## 安装

```bash
go get github.com/afeiship/go-yaml-path
```

## 快速开始

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

    // 创建 YPath 实例
    yp, err := ypath.NewFromString(yamlData)
    if err != nil {
        panic(err)
    }

    // 使用点号语法访问值
    fmt.Println("服务器主机:", yp.GetString("server.host"))        // localhost
    fmt.Println("服务器端口:", yp.GetInt("server.port"))           // 8080
    fmt.Println("SSL 启用:", yp.GetBool("server.ssl"))           // true
    fmt.Println("最大连接池:", yp.GetInt("database.connection.pool.max")) // 10

    // 数组访问
    fmt.Println("第一个功能:", yp.GetString("features.0"))       // authentication

    // 检查路径是否存在
    if yp.Exists("server.host") {
        fmt.Println("服务器主机已配置")
    }

    // 通配符访问
    features := yp.GetAll("features.*")
    fmt.Println("所有功能:", features)                          // [authentication logging monitoring]
}
```

## API 参考

### 创建 YPath 实例

```go
// 从字节数组创建
yp, err := ypath.New(yamlBytes)

// 从字符串创建
yp, err := ypath.NewFromString(yamlString)

// 从文件创建
yp, err := ypath.NewFromFile("config.yaml")
```

#### `NewFromFile(filename string) (*YPath, error)`
从 YAML 文件创建 YPath 实例。

```go
yp, err := ypath.NewFromFile("config.yaml")
if err != nil {
    panic(err)
}
```

### 基础方法

#### `Get(path string) interface{}`
检索指定路径的原始值。

```go
value := yp.Get("server.host")
```

#### `GetString(path string) string`
检索字符串值，转换其他类型为字符串。

```go
host := yp.GetString("server.host")  // "localhost"
port := yp.GetString("server.port")  // "8080" (从 int 转换)
```

#### `GetInt(path string) int`
检索整数值，支持类型转换。

```go
port := yp.GetInt("server.port")           // 8080
maxPool := yp.GetInt("database.pool.max")  // 10
```

#### `GetBool(path string) bool`
检索布尔值，支持字符串转换。

```go
ssl := yp.GetBool("server.ssl")      // true
active := yp.GetBool("server.active") // false
```

#### `GetFloat64(path string) float64`
检索 float64 值，支持类型转换。

```go
timeout := yp.GetFloat64("connection.timeout")  // 30.5
```

#### `Exists(path string) bool`
检查路径是否存在于 YAML 数据中。

```go
if yp.Exists("server.host") {
    // 路径存在
}
```

### 高级方法

#### `GetAll(path string) []interface{}`
检索匹配通配符路径的所有值。

```go
// 获取所有服务器主机
hosts := yp.GetAll("servers.*.host")
// 结果: []interface{}{"server1.example.com", "server2.example.com"}

// 获取所有功能
features := yp.GetAll("features.*")
// 结果: []interface{}{"authentication", "logging", "monitoring"}
```

## 路径语法

### 点号表示法
使用点号访问嵌套值：
- `server.host` → 访问 `server` 中的 `host`
- `database.connection.pool.max` → 深度嵌套访问

### 数组索引
使用数字索引访问数组元素：
- `features.0` → `features` 数组的第一个元素
- `servers.1.host` → 第二台服务器的 `host`

### 通配符
使用 `*` 作为通配符匹配多个值：
- `servers.*.host` → 所有服务器主机
- `features.*` → 所有功能
- `database.*.max` → `database` 下所有 `max` 值

### 特殊情况
- 空字符串 `""` 或点 `"."` → 返回根 YAML 数据
- 不存在的路径 → 返回适当的零值（nil, "", 0, false）

## 类型转换

库在可能的情况下自动转换类型：

- **GetString**: 任何类型 → 使用 `fmt.Sprintf` 转为字符串
- **GetInt**: int, int64, float64 → int；字符串 → 尝试解析为 int
- **GetBool**: bool → bool；字符串 → "true"/"1" 为 true；数字 → 非零为 true
- **GetFloat64**: float64, int, int64 → float64；字符串 → 尝试解析为 float64

## 错误处理

- **无效 YAML**: 返回 `yaml.Unmarshal` 错误
- **不存在的路径**: 返回零值，无错误
- **类型转换失败**: 返回零值，无错误

## 使用场景示例

### 配置文件
```go
yp, err := ypath.NewFromFile("config.yaml")
if err != nil {
    panic(err)
}
dbHost := yp.GetString("database.host")
dbPort := yp.GetInt("database.port")
enableSSL := yp.GetBool("server.ssl")
```

### Kubernetes 资源
```go
yp, _ := ypath.NewFromString(k8sYaml)
replicas := yp.GetInt("spec.replicas")
imageName := yp.GetString("spec.containers.0.image")
```

### API 响应
```go
yp, _ := ypath.NewFromString(apiResponse)
statusCode := yp.GetInt("response.status")
message := yp.GetString("response.message")
errors := yp.GetAll("response.errors.*.code")
```

## 运行示例

```bash
cd example
go run main.go
```

## 运行测试

```bash
go test -v
```

## 许可证

MIT 许可证
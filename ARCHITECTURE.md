# 文件架构说明

## 文件结构

```
go-yaml-path/
├── 源码文件 (根目录)
│   ├── ypath.go           # 核心定义和构造函数 (36 行)
│   ├── ypath_core.go      # 核心获取方法：Get, GetString, GetInt, GetBool, GetFloat64, Exists (128 行)
│   ├── ypath_wildcard.go  # 通配符功能：GetAll, collectWildcardValues (68 行)
│   └── ypath_list.go      # 列表方法：GetStringList, GetIntList, GetBoolList, GetFloat64List, GetList (228 行)
├── tests/                 # 测试目录
│   ├── test_utils.go      # 测试工具和常量
│   ├── construction_test.go    # 构造函数测试
│   ├── basic_getters_test.go    # 基础获取方法测试
│   ├── array_access_test.go     # 数组访问测试
│   ├── wildcard_test.go         # 通配符功能测试
│   ├── list_getters_test.go     # 列表方法测试
│   └── edge_cases_test.go       # 边界情况测试
├── example/               # 示例和演示
│   ├── main.go            # 使用示例
│   └── config.yaml        # 示例配置文件
├── docs/                  # 文档目录
│   └── 01-prd.md         # 产品需求文档
└── README*.md             # 项目文档
```

## 模块职责

### ypath.go - 核心定义
- `YPath` 结构体定义
- 构造函数：`New()`, `NewFromString()`, `NewFromFile()`
- 基础依赖：`gopkg.in/yaml.v3` 和 `os` 包

### ypath_core.go - 基础访问方法
- `Get()` - 通用路径访问方法
- `GetString()` - 字符串值获取
- `GetInt()` - 整数值获取
- `GetBool()` - 布尔值获取
- `GetFloat64()` - 浮点数获取
- `Exists()` - 路径存在检查
- 路径解析和导航逻辑

### ypath_wildcard.go - 通配符支持
- `GetAll()` - 通配符路径匹配
- `collectWildcardValues()` - 递归通配符值收集
- `*` 通配符语法支持
- 数组和对象通配符遍历

### ypath_list.go - 列表类型方法
- `GetStringList()` - 字符串列表获取，支持通配符
- `GetIntList()` - 整数列表获取，支持类型转换和通配符
- `GetBoolList()` - 布尔列表获取，支持类型转换和通配符
- `GetFloat64List()` - 浮点数列表获取，支持类型转换和通配符
- `GetList()` - 通用列表方法，支持通配符
- 单值到列表的自动转换

## 测试架构

### 测试文件组织
- **test_utils.go**: 测试工具函数、常量和辅助方法
- **construction_test.go**: 构造函数相关测试（New, NewFromString, NewFromFile）
- **basic_getters_test.go**: 基础获取方法测试（GetString, GetInt, GetBool, GetFloat64, Get, Exists）
- **array_access_test.go**: 数组和索引访问测试
- **wildcard_test.go**: 通配符功能测试（GetAll）
- **list_getters_test.go**: 列表方法测试（GetStringList, GetIntList, GetBoolList, GetFloat64List, GetList）
- **edge_cases_test.go**: 边界情况和异常处理测试

### 测试覆盖
- ✅ 23 个测试函数，覆盖所有主要功能
- ✅ 正常情况、异常情况、边界情况全覆盖
- ✅ 临时文件创建和清理
- ✅ 复杂嵌套结构测试

## 设计优势

1. **模块化**: 每个文件专注于特定功能域
2. **可维护性**: 代码组织清晰，易于理解和修改
3. **可扩展性**: 新功能可以独立添加到相应模块
4. **测试性**: 测试文件按功能拆分，便于定位和维护
5. **性能**: 减少单个文件大小，提高编译和测试效率

## 使用示例

```go
// 创建实例 (ypath.go)
yp, err := ypath.NewFromFile("config.yaml")

// 基础访问 (ypath_core.go)
host := yp.GetString("server.host")
port := yp.GetInt("server.port")

// 通配符访问 (ypath_wildcard.go)
allPorts := yp.GetAll("servers.*.port")

// 列表访问 (ypath_list.go)
portList := yp.GetIntList("servers.*.port")
features := yp.GetStringList("features")
```
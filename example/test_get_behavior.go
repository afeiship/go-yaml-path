package main

import (
	"fmt"
	"log"
	"strings"

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
features:
  - authentication
  - logging
`

	yp, err := ypath.NewFromString(yamlData)
	if err != nil {
		log.Fatalf("Failed to create YPath: %v", err)
	}

	fmt.Println("=== yp.Get() 方法行为测试 ===")
	fmt.Println()

	// 1. 传入空字符串 - 返回整个 YAML 数据的根对象
	fmt.Println("1. yp.Get(\"\"):")
	rootData := yp.Get("")
	fmt.Printf("   类型: %T\n", rootData)
	fmt.Printf("   值: %+v\n", rootData)
	fmt.Println()

	// 2. 传入点号 - 同样返回整个 YAML 数据的根对象
	fmt.Println("2. yp.Get(\".\"):")
	dotData := yp.Get(".")
	fmt.Printf("   类型: %T\n", dotData)
	fmt.Printf("   值: %+v\n", dotData)
	fmt.Println()

	// 3. 传入具体路径 - 返回对应的值
	fmt.Println("3. yp.Get(\"server.host\"):")
	host := yp.Get("server.host")
	fmt.Printf("   类型: %T\n", host)
	fmt.Printf("   值: %v\n", host)
	fmt.Println()

	fmt.Println("4. yp.Get(\"features\"):")
	features := yp.Get("features")
	fmt.Printf("   类型: %T\n", features)
	fmt.Printf("   值: %v\n", features)
	fmt.Println()

	// 4. 传入不存在的路径 - 返回 nil
	fmt.Println("5. yp.Get(\"nonexistent.path\"):")
	nonExistent := yp.Get("nonexistent.path")
	fmt.Printf("   类型: %T\n", nonExistent)
	fmt.Printf("   值: %v\n", nonExistent)
	fmt.Println()

	fmt.Println("=== 总结 ===")
	fmt.Println("• yp.Get(\"\") - 返回整个 YAML 根对象")
	fmt.Println("• yp.Get(\".\") - 返回整个 YAML 根对象")
	fmt.Println("• yp.Get(\"path\") - 返回指定路径的值")
	fmt.Println("• yp.Get(\"nonexistent\") - 返回 nil")
	fmt.Println("• yp.Get() 必须传入参数，因为函数签名需要 string 参数")
	fmt.Println()
	fmt.Println("=== Get() 方法源码分析 ===")
	fmt.Println(strings.Repeat(" ", 4) + "if path == \"\" || path == \".\" {")
	fmt.Println(strings.Repeat(" ", 8) + "return yp.data")
	fmt.Println(strings.Repeat(" ", 4) + "}")
	fmt.Println()
	fmt.Println("当传入空字符串或点号时，直接返回 yp.data（整个 YAML 的解析结果）")
}
package main

import (
	"fmt"
	"log"

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
`

	yp, err := ypath.NewFromString(yamlData)
	if err != nil {
		log.Fatalf("Failed to create YPath: %v", err)
	}

	fmt.Println("=== Get() 无参数方法演示 ===")
	fmt.Println()

	// 演示 Get() 无参数方法
	fmt.Println("1. yp.Get() - 无参数方法:")
	rootData := yp.Get()
	fmt.Printf("   类型: %T\n", rootData)
	fmt.Printf("   是否为空: %t\n", rootData == nil)
	fmt.Println()

	fmt.Println("2. yp.Get(\"\") - 空字符串方法:")
	emptyPathData := yp.Get("")
	fmt.Printf("   类型: %T\n", emptyPathData)
	fmt.Printf("   是否为空: %t\n", emptyPathData == nil)
	fmt.Println()

	fmt.Println("3. yp.Get(\".\") - 点号方法:")
	dotPathData := yp.Get(".")
	fmt.Printf("   类型: %T\n", dotPathData)
	fmt.Printf("   是否为空: %t\n", dotPathData == nil)
	fmt.Println()

	fmt.Println("4. 三种方式的一致性:")
	dataMap, _ := rootData.(map[string]interface{})
	emptyMap, _ := emptyPathData.(map[string]interface{})
	dotMap, _ := dotPathData.(map[string]interface{})

	sameSize := len(dataMap) == len(emptyMap) && len(dataMap) == len(dotMap)
	fmt.Printf("   三种方式返回相同大小的数据: %t\n", sameSize)
	fmt.Println()

	fmt.Println("5. 使用 Get() 无参数的优势:")
	fmt.Printf("   yp.Get() 获取根对象，然后进一步操作: %v\n", rootData)
	fmt.Println("   语法最简洁，无需任何参数")
	fmt.Println("   表达意图最清晰")
	fmt.Println("   与传统 Get(\"\") 完全等价")
	fmt.Println()

	fmt.Println("6. 实际使用场景:")
	fmt.Println("   // 获取完整配置")
	fmt.Println("   config := yp.Get()")
	fmt.Println("   ")
	fmt.Println("   // 转换为 JSON")
	fmt.Println("   jsonBytes, _ := json.Marshal(yp.Get())")
	fmt.Println("   ")
	fmt.Println("   // 调试打印")
	fmt.Println("   fmt.Printf(\"Full config: %+v\\n\", yp.Get())")
	fmt.Println("   ")
	fmt.Println("   // 传统方式仍然有效")
	fmt.Println("   config := yp.Get(\"\")  // 或 yp.Get(\".\")")
}
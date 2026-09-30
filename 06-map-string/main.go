/*
第 06 课：map 与字符串

map：键值对。写入前必须 make，否则 panic。
string：不可变；常用操作在 strings 包。

运行：
  go run ./06-map-string
*/

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== 第 06 课：map 与字符串 ===")
	demoMap()
	demoString()
}

func demoMap() {
	fmt.Println("\n--- map 增删改查 ---")

	// 必须 make（或用字面量），nil map 不能写
	scores := make(map[string]int)
	scores["张三"] = 90
	scores["李四"] = 80

	fmt.Println("张三 =", scores["张三"])

	// comma-ok：区分「键不存在」和「值就是 0」
	v, ok := scores["王五"]
	if !ok {
		fmt.Println("王五 不存在")
	} else {
		fmt.Println("王五 =", v)
	}

	scores["李四"] = 85 // 改
	delete(scores, "张三") // 删

	fmt.Println("遍历（顺序不保证）：")
	for name, score := range scores {
		fmt.Printf("  %s = %d\n", name, score)
	}
	fmt.Println("len =", len(scores))
}

func demoString() {
	fmt.Println("\n--- 字符串常用操作 ---")

	s := "Hello Go Go"
	fmt.Println("拼接:", "Hi, " + "Go")
	fmt.Println("Contains:", strings.Contains(s, "Go"))
	fmt.Println("Index:", strings.Index(s, "Go"))
	fmt.Println("Split:", strings.Split("a,b,c", ","))
	fmt.Println("Replace:", strings.ReplaceAll(s, "Go", "Golang"))
	fmt.Println("ToUpper:", strings.ToUpper("go"))
	fmt.Println("TrimSpace:", strings.TrimSpace("  hi  "))

	// range 字符串按 rune 走；中文一个字可能占多字节
	for i, ch := range "Go好" {
		fmt.Printf("  字节下标=%d 字符=%c\n", i, ch)
	}
}

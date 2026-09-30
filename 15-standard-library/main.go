/*
第 15 课：常用标准库

- strconv：字符串 ↔ 数字
- time：当前时间、格式化（模板死记 2006-01-02 15:04:05）
- os：读写文件
- encoding/json：结构体 ↔ JSON（Gin 绑 JSON 的底层就是这个）

运行：
  go run ./15-standard-library
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	fmt.Println("=== 第 15 课：标准库 ===")
	demoStrconv()
	demoTime()
	demoJSON()
	demoFile()
}

func demoStrconv() {
	fmt.Println("\n--- strconv ---")

	n, err := strconv.Atoi("123") // ASCII to int
	if err != nil {
		fmt.Println("Atoi 失败:", err)
	} else {
		fmt.Println("Atoi(\"123\") =", n)
	}

	_, err = strconv.Atoi("abc")
	fmt.Println("Atoi(\"abc\") 错误:", err)

	fmt.Println("Itoa(45) =", strconv.Itoa(45)) // int to ASCII
}

func demoTime() {
	fmt.Println("\n--- time ---")
	now := time.Now()
	fmt.Println("现在:", now.Format("2006-01-02 15:04:05"))
	fmt.Println("时间戳(秒):", now.Unix())
}

type userDTO struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func demoJSON() {
	fmt.Println("\n--- encoding/json ---")
	u := userDTO{Name: "小明", Age: 18}
	b, err := json.Marshal(u) // 结构体 → JSON
	if err != nil {
		fmt.Println("Marshal 失败:", err)
		return
	}
	fmt.Println("结构体 → JSON:", string(b))

	var u2 userDTO
	if err := json.Unmarshal([]byte(`{"name":"小红","age":20}`), &u2); err != nil {
		fmt.Println("Unmarshal 失败:", err)
		return
	}
	fmt.Printf("JSON → 结构体: %+v\n", u2)
}

func demoFile() {
	fmt.Println("\n--- os 文件读写 ---")
	path := filepath.Join(os.TempDir(), "go-study-15-demo.txt")
	if err := os.WriteFile(path, []byte("hello go\n"), 0644); err != nil {
		fmt.Println("写入失败:", err)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}
	fmt.Printf("读到: %s", data)
	_ = os.Remove(path)
}

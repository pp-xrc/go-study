/*
第 08 课：变量作用域

- 包级变量：本文件（本包）都能用
- 函数内变量：出了函数就没了
- 短声明 := 在内层会「新造变量」还是「赋值」，看左边有没有新名字

运行：
  go run ./08-scope
*/

package main

import "fmt"

// 包级变量：main 和其他函数都能读
var packageName = "go-study"



func main() {
	fmt.Println("=== 第 08 课：作用域 ===")

	demoLocal()
	demoInnerBlock()
	demoShortDeclShadow()
}

func demoLocal() {
	fmt.Println("\n--- 函数内局部变量 ---")
	x := 10
	fmt.Println("包级 packageName =", packageName)
	fmt.Println("函数内 x =", x)
	printX() // 这里访问不到上面的 x
}

func printX() {
	// fmt.Println(x) // 编译错误：x 未定义
	fmt.Println("printX 里只能看到包级变量:", packageName)
}

func demoInnerBlock() {
	fmt.Println("\n--- if / for 的花括号也是一层作用域 ---")
	n := 1
	if n > 0 {
		msg := "n 是正数" // msg 只在 if 里面活
		fmt.Println(msg)
	}
	// fmt.Println(msg) // 编译错误
}

func demoShortDeclShadow() {
	fmt.Println("\n--- := 的「影子变量」---")
	err := fmt.Errorf("外层 err")
	if true {
		// 左边出现了新名字 n，所以 err 也是内层新变量，不是给外层赋值
		n, err := 1, fmt.Errorf("内层 err")
		fmt.Println("  内层 n, err =", n, err)
	}
	fmt.Println("  外层 err 仍是:", err)
}

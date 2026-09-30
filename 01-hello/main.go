/*
第 01 课：Hello World、常量、变量、条件判断

运行：
  go run ./01-hello
*/

package main

import "fmt" // 标准库：打印输出

// const 常量：编译期就定死，程序运行中不能改
const a = 10
const b = 20
const sum = a + b // 常量之间可以互相运算

func main() {
	fmt.Println("=== 第 01 课：Hello World ===")

	// Println：打印后换行
	// Printf：按占位符格式化。%d 整数，%s 字符串，%T 类型
	fmt.Printf("Hello, World! sum = %d\n", sum)

	// if / else：条件必须是 bool，括号可以省略
	if sum > 10 {
		fmt.Println("sum 大于 10")
	} else {
		fmt.Println("sum 不大于 10")
	}

	// var 变量：运行时才能确定，可以改
	var score = 10
	var name = "小明"
	score = 95 // 变量可以重新赋值；常量不行

	// Sprintf：不打印，只把格式化结果返回成字符串
	info := fmt.Sprintf("姓名: %s, 分数: %d", name, score)
	fmt.Println(info)
}

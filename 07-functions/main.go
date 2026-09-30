/*
第 07 课：函数

要点：Go 全是值传递；多返回值处理错误；函数可以当「值」传来传去。

运行：
  go run ./07-functions
*/

package main

import "fmt"

func main() {
	fmt.Println("=== 第 07 课：函数 ===")

	demoBasic()
	demoVariadic()
	demoFuncValue()
	demoAnonymous()
	demoDefer()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 可变参数：函数内部 nums 就是 []int
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func demoBasic() {
	fmt.Println("\n--- 普通函数 ---")
	fmt.Println("max(3, 9) =", max(3, 9))
}

func demoVariadic() {
	fmt.Println("\n--- 可变参数 ...int ---")
	fmt.Println("sum(1,2,3) =", sum(1, 2, 3))

	arr := []int{10, 20, 30}
	fmt.Println("sum(arr...) =", sum(arr...)) // 切片打散传入
}

func demoFuncValue() {
	fmt.Println("\n--- 函数当作参数 ---")
	apply := func(f func(int, int) int, x, y int) int {
		return f(x, y)
	}
	fmt.Println("apply(max, 4, 7) =", apply(maxInt, 4, 7))
}

func demoAnonymous() {
	fmt.Println("\n--- 匿名函数 ---")
	add := func(a, b int) int {
		return a + b
	}
	fmt.Println("add(2, 3) =", add(2, 3))
}

func demoDefer() {
	fmt.Println("\n--- defer：函数结束前执行，多个 defer 后进先出 ---")
	defer fmt.Println("  defer 1（最后打印）")
	defer fmt.Println("  defer 2")
	fmt.Println("  函数体先执行")
}

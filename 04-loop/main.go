/*
第 04 课：for 循环与 range

Go 只有 for，没有 while。range 用来遍历切片、字符串、map；Go 1.22+ 也可用 for range n 重复 n 次。

运行：
  go run ./04-loop
*/

package main

import "fmt"

func main() {
	fmt.Println("=== 第 04 课：循环 ===")

	demoForClassic()
	demoForWhile()
	demoRangeSlice()
	demoRangeString()
	demoRangeInt()
	demoBreakContinue()
}

func demoForClassic() {
	fmt.Println("\n--- 经典 for：初始值; 条件; 步进 ---")
	sum := 0
	for i := 1; i <= 5; i++ {
		sum += i
	}
	fmt.Println("1+2+3+4+5 =", sum)
}

func demoForWhile() {
	fmt.Println("\n--- 只写条件：相当于 while ---")
	n := 3
	for n > 0 {
		fmt.Println("n =", n)
		n--
	}
}

func demoRangeSlice() {
	fmt.Println("\n--- range 切片：得到 (下标, 值) ---")
	nums := []int{10, 20, 30}
	for i, v := range nums {
		fmt.Printf("下标=%d 值=%d\n", i, v)
	}
	// 只要值、不要下标：写成 for _, v := range nums
}

func demoRangeString() {
	fmt.Println("\n--- range 字符串：按字符（rune）遍历，不是按字节 ---")
	s := "Go语言"
	for i, ch := range s {
		fmt.Printf("字节下标=%d 字符=%c\n", i, ch)
	}
}

func demoRangeInt() {
	fmt.Println("\n--- Go 1.22+：for range n 只重复 n 次（不用下标时）---")
	for range 3 {
		fmt.Println("  ping")
	}
}

func demoBreakContinue() {
	fmt.Println("\n--- break 跳出循环；continue 跳过本轮 ---")
	for i := 1; i <= 5; i++ {
		if i == 2 {
			continue // 不打印 2
		}
		if i == 4 {
			break // 到 4 就结束，5 不会跑
		}
		fmt.Println("i =", i)
	}
}

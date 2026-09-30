/*
第 05 课：数组与切片

数组：长度写死在类型里，[3]int 和 [4]int 是不同类型。
切片：长度可变，日常开发几乎都用切片。

运行：
  go run ./05-arr-slice
*/

package main

import "fmt"

func main() {
	fmt.Println("=== 第 05 课：数组与切片 ===")

	demoArray()
	demoSlice()
	demoAppendLenCap()
	demoSliceFromArray()
}

func demoArray() {
	fmt.Println("\n--- 数组：定长 ---")
	arr := [3]int{1, 2, 3}
	arr[0] = 10
	fmt.Println("arr =", arr, "len =", len(arr))

	for i, v := range arr {
		fmt.Printf("  [%d]=%d\n", i, v)
	}
}

func demoSlice() {
	fmt.Println("\n--- 切片：变长，三种创建方式 ---")

	s1 := []int{1, 2, 3}    // 字面量
	s2 := make([]int, 2, 5) // make(类型, 长度, 容量capacity)  长度个元素已是零值
	fmt.Println("字面量 s1 =", s1)
	fmt.Println("make   s2 =", s2, "len=", len(s2), "cap=", cap(s2))
}

func demoAppendLenCap() {
	fmt.Println("\n--- append / 截取 / copy ---")

	s := []int{1, 2, 3}
	s = append(s, 4, 5) // 追加后必须接住返回值
	fmt.Println("append 后", s)

	part := s[1:4] // 左闭右开：下标 1,2,3 → 2,3,4
	fmt.Println("s[1:4] =", part)

	dst := make([]int, len(part))
	copy(dst, part) // 拷贝到新切片，互不影响
	fmt.Println("copy 后 dst =", dst)
	fmt.Println("len =", len(s), "cap =", cap(s))
}

func demoSliceFromArray() {
	fmt.Println("\n--- 从数组切出切片（共享底层数组，改切片可能改到数组）---")
	arr := [4]int{10, 20, 30, 40}
	s := arr[0:2] // [10 20]
	s[0] = 99
	fmt.Println("切片 s =", s)
	fmt.Println("数组 arr =", arr) // arr[0] 也变成 99
}

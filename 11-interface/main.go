/*
第 11 课：接口

接口 = 方法集合。谁实现了这些方法，谁就是这个接口。
不必写 implements，编译器自动认。

运行：
  go run ./11-interface
*/

package main

import "fmt"

type Animal interface {
	Sleep()
	Speak() string
}

type Cat struct{ Name string }

func (c Cat) Sleep()        { fmt.Println(c.Name, "在睡觉") }
func (c Cat) Speak() string { return "喵" }

type Dog struct{ Name string }

func (d Dog) Sleep()        { fmt.Println(d.Name, "在睡觉") }
func (d Dog) Speak() string { return "汪" }

// 参数是接口：传入 Cat 或 Dog 都行（多态）
func play(a Animal) {
	a.Sleep()
	fmt.Println("叫声:", a.Speak())
}

func printAny(v any) {
	// comma-ok 断言：失败不会 panic
	if s, ok := v.(string); ok {
		fmt.Println("这是字符串:", s)
		return
	}

	switch x := v.(type) {
	case int:
		fmt.Println("这是 int:", x)
	case Animal:
		fmt.Println("这是 Animal，叫声:", x.Speak())
	default:
		fmt.Printf("其他类型 %T: %v\n", v, v)
	}
}

func main() {
	fmt.Println("=== 第 11 课：接口 ===")

	fmt.Println("\n--- 多态 ---")
	play(Cat{Name: "小花"})
	play(Dog{Name: "旺财"})

	fmt.Println("\n--- 空接口 any：什么值都能接 ---")
	printAny("hello")
	printAny(42)
	printAny(Cat{Name: "小花"})
	printAny([]int{1, 2})
}

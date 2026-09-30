/*
指针课（配合第 07 课函数、第 09 课结构体）

指针 = 一个变量，里面存的是「另一个变量的内存地址」。

常用符号：
  &x   取地址：得到 x 在内存里的位置
  *T   指针类型：指向 T 的指针，例如 *int
  *p   解引用：顺着地址找到原来的值；也可以给这个位置赋值

建议学完第 07 课（函数传参）和第 09 课（指针接收者）后再对照着看。

运行：
  go run ./07b-pointers
*/

package main

import "fmt"

type Student struct {
	Name  string
	Score int
}

func main() {
	fmt.Println("=== 指针 ===")
	demoBasic()
	demoNil()
	demoPass()
	demoStruct()
}

func demoBasic() {
	fmt.Println("\n--- 1. 取址、解引用、通过指针修改 ---")

	a := 10
	var p *int = &a // p 的类型是 *int，值是 a 的地址

	fmt.Printf("a 的值     = %d\n", a)
	fmt.Printf("a 的地址   = %p\n", &a)
	fmt.Printf("p 存的地址 = %p  （应和 &a 相同）\n", p)
	fmt.Printf("*p 解引用  = %d  （顺着地址取出 a）\n", *p)

	*p = 20 // 改的是 a 那块内存
	fmt.Println("执行 *p = 20 之后，a =", a)
}

func demoNil() {
	fmt.Println("\n--- 2. 指针零值是 nil，解引用会 panic，先判断 ---")

	var p *int
	fmt.Println("未赋值的指针 p == nil ?", p == nil)
	if p!= nil {
		fmt.Println(*p)
	} else {
		fmt.Println("p 是 nil，不能写 *p")
	}
}

func addOneByValue(n int) {
	n++ // n 是副本，外面的变量不变
}

func addOneByPointer(n *int) {
	*n++ // 改原来的那份
}

func demoPass() {
	fmt.Println("\n--- 3. 函数传参：值 vs 指针 ---")
	fmt.Println("Go 全是值传递；想改外面的变量，就传地址")

	x := 1
	addOneByValue(x)
	fmt.Println("addOneByValue(x)  之后 x =", x) // 仍是 1

	addOneByPointer(&x) // 传入 x 的地址
	fmt.Println("addOneByPointer(&x) 之后 x =", x) // 变成 2
}

func (s Student) setScoreByValue(score int) {
	s.Score = score // 改副本
}

func (s *Student) setScoreByPointer(score int) {
	s.Score = score // 改原来的结构体
}

func demoStruct() {
	fmt.Println("\n--- 4. 结构体：方法要用指针接收者才能改字段 ---")

	s := Student{Name: "小明", Score: 60}
	s.setScoreByValue(100)
	fmt.Println("值接收者改分后  Score =", s.Score) // 仍是 60

	s.setScoreByPointer(100)
	fmt.Println("指针接收者改分后 Score =", s.Score) // 100

	// 访问指针上的字段不必写 (*p).Name，Go 会自动解引用
	p := &s
	fmt.Println("p.Name =", p.Name)
}

/*
第 02 课：基本类型、短声明、零值、iota

运行：
  go run ./02-basic-types
*/

package main

import "fmt"

func main() {
	fmt.Println("=== 第 02 课：基本数据类型 ===")

	demoTypes()
	demoZeroValue()
	demoShortDecl()
	demoIota()
}

func demoTypes() {
	fmt.Println("\n--- 常见类型 ---")

	var n int = 10          // 整数
	var f float64 = 3.14    // 浮点数（默认推荐 float64）
	var ok bool = true      // 布尔
	var s string = "hello"  // 字符串
	var ch rune = '中'       // 字符，底层是 int32（Unicode 码点）
	var bs []byte = []byte(s) // 字节切片，常用来处理二进制 / UTF-8

	fmt.Printf("int     n=%d   类型=%T\n", n, n)
	fmt.Printf("float64 f=%v   类型=%T\n", f, f)
	fmt.Printf("bool    ok=%t  类型=%T\n", ok, ok)
	fmt.Printf("string  s=%s   类型=%T\n", s, s)
	fmt.Printf("rune    ch=%c  类型=%T  码点=%d\n", ch, ch, ch)
	fmt.Printf("[]byte  bs=%v  类型=%T\n", bs, bs)
}

func demoZeroValue() {
	fmt.Println("\n--- 零值：声明了但不赋值时的默认值 ---")

	var n int
	var f float64
	var ok bool
	var s string
	fmt.Printf("int=%d  float64=%v  bool=%t  string=%q\n", n, f, ok, s)
}

func demoShortDecl() {
	fmt.Println("\n--- 短声明 := （只能写在函数里）---")

	x, y := 10, 20 // 类型由右边自动推断
	fmt.Printf("x=%d (%T)  y=%d (%T)\n", x, x, y, y)

	x, y = y, x // 交换：两边类型必须相同
	fmt.Printf("交换后 x=%d y=%d\n", x, y)
}

func demoIota() {
	fmt.Println("\n--- iota：const 块里的计数器，每行 +1，新的 const 从 0 开始 ---")

	const (
		Sun = iota // 0
		Mon        // 1  省略右边时，沿用上一行公式
		Tue        // 2
	)
	fmt.Println("Sun, Mon, Tue =", Sun, Mon, Tue)

	const (
		Sunday = iota + 1 // 从 1 开始：1,2,3...
		Monday
		Tuesday
	)
	fmt.Println("Sunday, Monday, Tuesday =", Sunday, Monday, Tuesday)
}

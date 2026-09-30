/*
第 12 课：包与模块

- 每个目录一个包；package main + func main() 才是可执行程序
- 导入路径 = go.mod 里的模块名 + 相对目录
- 首字母大写才能被别的包用（导出）

运行：
  go run ./12-package-module
*/

package main

import (
	"fmt"

	"go-study/12-package-module/calc"
	hello "go-study/12-package-module/greeter" // 别名：用 hello 代替 greeter

	"rsc.io/quote" // 第三方包，版本锁在 go.mod
)

func main() {
	fmt.Println("=== 第 12 课：包与模块 ===")

	fmt.Println("\n--- 自定义包 calc ---")
	fmt.Println("Version =", calc.Version)
	fmt.Println("Add(3,5) =", calc.Add(3, 5))
	fmt.Println("Mul(3,4) =", calc.Mul(3, 4)) // 内部调用了未导出的 multiply

	fmt.Println("\n--- 自定义包 greeter（别名 hello）---")
	p := hello.NewPerson("小明", 18)
	fmt.Println("Name =", p.Name)   // 大写字段，包外能直接读
	fmt.Println("Age() =", p.Age()) // 小写字段只能通过方法读
	fmt.Println(p.Introduce())
	// p.age = 20 // 编译错误

	fmt.Println("\n--- 第三方包 ---")
	fmt.Println(quote.Hello())
	// calc.multiply(2, 3) // 编译错误：小写，包外不可见

}

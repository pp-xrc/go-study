/*
第 03 课：运算符、switch、多返回值 error

运行：
  go run ./03-operators
*/

package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("=== 第 03 课：运算符 ===")

	demoCompare()
	demoBit()
	demoSwitch()
	demoCalculator()
}

func demoCompare() {
	fmt.Println("\n--- 比较运算 ---")
	a, b := 10, 20
	fmt.Println("a == b", a == b)
	fmt.Println("a != b", a != b)
	fmt.Println("a <  b", a < b)
	fmt.Println("a >= b", a >= b)
}

func demoBit() {
	fmt.Println("\n--- 位运算（以 10=1010, 20=10100 为例，心里有个数即可）---")
	a, b := 10, 20
	fmt.Println("a & b  按位与   =", a&b)
	fmt.Println("a | b  按位或   =", a|b)
	fmt.Println("a ^ b  按位异或 =", a^b)
	fmt.Println("a << 2 左移 2 位(≈×4) =", a<<2)
	fmt.Println("a >> 1 右移 1 位(≈÷2) =", a>>1)
}

func demoSwitch() {
	fmt.Println("\n--- switch：不用写 break，匹配到就结束 ---")
	score := 85
	switch {
	case score >= 90:
		fmt.Println("优秀")
	case score >= 60:
		fmt.Println("及格")
	default:
		fmt.Println("不及格")
	}
}

// calculate 返回 (结果, 错误)。成功时 error 是 nil
func calculate(a, b float64, op rune) (float64, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, fmt.Errorf("除数不能为 0")
		}
		return a / b, nil
	case '%':
		if b == 0 {
			return 0, fmt.Errorf("除数不能为 0")
		}
		return math.Mod(a, b), nil
	default:
		return 0, fmt.Errorf("不支持的运算符: %c", op)
	}
}

func demoCalculator() {
	fmt.Println("\n--- 多返回值：(结果, error)，先判断 err ---")

	cases := []struct {
		a, b float64
		op   rune
	}{
		{10, 3, '+'},
		{10, 0, '/'},
		{10, 3, '?'},
	}
	for _, c := range cases {
		result, err := calculate(c.a, c.b, c.op)
		if err != nil {
			fmt.Printf("%.0f %c %.0f → 失败: %v\n", c.a, c.op, c.b, err)
			continue
		}
		fmt.Printf("%.0f %c %.0f = %.2f\n", c.a, c.op, c.b, result)
	}
}

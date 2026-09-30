// Package calc 自定义包示例：加减乘
package calc

// Version 大写 → 导出，其他包能读
const Version = "1.0.0"

// Add 导出函数
func Add(a, b int) int { return a + b }

func Sub(a, b int) int { return a - b }

// multiply 未导出，包外不能直接调
func multiply(a, b int) int { return a * b }

// Mul 导出，内部可以调用未导出函数
func Mul(a, b int) int { return multiply(a, b) }

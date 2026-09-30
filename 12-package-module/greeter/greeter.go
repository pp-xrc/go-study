// Package greeter 演示：导出类型、导出/未导出字段
package greeter

import "fmt"

// Person 导出结构体
type Person struct {
	Name string // 大写：包外可读可写
	age  int    // 小写：包外不能直接碰
}

// NewPerson 构造函数：负责给未导出字段赋值
func NewPerson(name string, age int) Person {
	return Person{Name: name, age: age}
}

func (p Person) Age() int { return p.age }

func (p Person) Introduce() string {
	return fmt.Sprintf("你好，我是 %s，今年 %d 岁", p.Name, p.age)
}

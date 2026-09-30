/*
第 09 课：结构体与方法

结构体 = 自定义的一组字段。
方法挂在类型上：func (接收者 类型) 方法名()
要改字段必须用指针接收者，否则改的是副本。

运行：
  go run ./09-struct
*/

package main

import "fmt"

type Address struct {
	Province string
	City     string
}

type Student struct {
	Name    string
	Age     int
	Score   int
	Address Address // 嵌套：学生「包含」地址
}

// 值接收者：读数据。s 是副本，改了不影响外面
func (s Student) GetScore() int {
	return s.Score
}

// 指针接收者：改数据。s 指向原来那份
func (s *Student) SetScore(score int) {
	s.Score = score
}

func (s Student) Introduce() {
	fmt.Printf("我是%s，%d岁，分数%d，家在%s%s\n",
		s.Name, s.Age, s.Score, s.Address.Province, s.Address.City)
}

func main() {
	fmt.Println("=== 第 09 课：结构体 ===")

	// 1) 零值：字段都是各自类型的零值
	var s1 Student
	fmt.Printf("零值: %+v\n", s1)

	// 2) 字面量
	s2 := Student{
		Name:  "小明",
		Age:   18,
		Score: 90,
		Address: Address{
			Province: "浙江",
			City:     "杭州",
		},
	}
	s2.Introduce()

	// 3) 指针接收者才能改成功
	s2.SetScore(100)
	fmt.Println("SetScore 后:", s2.GetScore())

	// 4) new 得到 *Student，用起来和值差不多（Go 会自动解引用）
	s3 := new(Student)
	s3.Name = "小红"
	fmt.Println("new 出来的:", s3.Name)
	s3.Score = 95
	fmt.Println("new 出来的:", s3.Score)
	s3.Introduce()
}

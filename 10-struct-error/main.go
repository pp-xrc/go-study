/*
第 10 课：结构体 + map + error 综合

用 map 当「内存数据库」，做学生的增删查。
查不到、分数非法时返回 error。
末尾用 recover 接住 nil map 的 panic，进程不会退出。

运行：
  go run ./10-struct-error
*/

package main

import (
	"errors"
	"fmt"
)

type Student struct {
	ID    int
	Name  string
	Age   int
	Score int
}

type StudentManager struct {
	students map[int]Student
	nextID   int
}

// 构造函数：map 必须 make，否则写入会 panic
func NewStudentManager() *StudentManager {
	return &StudentManager{
		students: make(map[int]Student),
		nextID:   1,
	}
}

func (m *StudentManager) Add(name string, age, score int) Student {
	s := Student{ID: m.nextID, Name: name, Age: age, Score: score}
	m.students[s.ID] = s
	m.nextID++
	return s
}

func (m *StudentManager) Get(id int) (Student, error) {
	s, ok := m.students[id]
	if !ok {
		return Student{}, fmt.Errorf("学生 id=%d 不存在", id)
	}
	return s, nil
}

func (m *StudentManager) Delete(id int) error {
	if _, ok := m.students[id]; !ok {
		return fmt.Errorf("学生 id=%d 不存在", id)
	}
	delete(m.students, id)
	return nil
}

func (m *StudentManager) List() []Student {
	list := make([]Student, 0, len(m.students))
	for _, s := range m.students {
		list = append(list, s)
	}
	return list
}

func (s *Student) SetScore(score int) error {
	if score < 0 || score > 100 {
		return errors.New("分数必须在 0~100")
	}
	s.Score = score
	return nil
}

func main() {
	fmt.Println("=== 第 10 课：学生管理（map + error）===")

	mgr := NewStudentManager()
	s1 := mgr.Add("张三", 18, 90)
	mgr.Add("李四", 19, 85)

	if s, err := mgr.Get(s1.ID); err != nil {
		fmt.Println("查询失败:", err)
	} else {
		fmt.Printf("查到: %+v\n", s)
	}

	if _, err := mgr.Get(999); err != nil {
		fmt.Println("查询失败:", err)
	}

	fmt.Println("全部学生:")
	for _, s := range mgr.List() {
		fmt.Printf("  id=%d name=%s score=%d\n", s.ID, s.Name, s.Score)
	}

	_ = mgr.Delete(s1.ID)
	fmt.Println("删除后人数:", len(mgr.List()))

	st := Student{Name: "测试", Score: 80}
	if err := st.SetScore(-1); err != nil {
		fmt.Println("改分失败:", err)
	}

	demoRecoverNilMap()
}

// 写入 nil map 会 panic；defer + recover 能接住，后面的程序还能继续
func demoRecoverNilMap() {
	fmt.Println("\n--- panic / recover ---")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("抓住 panic:", r)
			fmt.Println("进程没有退出，可以继续跑")
		}
	}()
	var m map[int]Student
	m[1] = Student{ID: 1, Name: "会崩"}
	fmt.Println("这行打印不到")
}

/*
第 25 课：GORM 进阶 —— 条件 / 分页 / 排序 / 计数

================================================================================
一、学习内容
================================================================================
  - Where() / Or() 条件查询（多个 Where = AND）
  - Limit() / Offset() 分页（对接前端 page、size）
  - Order() 排序
  - Count() 统计总数

================================================================================
二、实操任务
================================================================================
  1. 条件：性别=女 且 age>18；姓名包含「张」
  2. 分页列表：接收 page、size
  3. 按 create_time 倒序，并返回 total

================================================================================
三、运行
================================================================================
  go run ./25-GORM-advance

================================================================================
*/

package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 映射 user 表（本课新增 gender 字段）
type User struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string    `gorm:"column:name;size:64;not null"`
	Age        uint      `gorm:"column:age;not null;default:0"`
	Gender     string    `gorm:"column:gender;size:8;not null;default:''"` // 男 / 女
	Phone      string    `gorm:"column:phone;size:20;not null;default:''"`
	Password   string    `gorm:"column:password;size:255;not null"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
}

func (User) TableName() string { return "user" }

// PageResult 分页返回：列表 + 总数
type PageResult struct {
	List  []User `json:"list"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

func main() {
	db := mustConnect()
	prepareSchemaAndDemoData(db)

	fmt.Println("========== 1. 条件查询：女 且 age > 18 ==========")
	females := findAdultFemales(db)
	printUsers(females)

	fmt.Println("========== 2. 条件查询：姓名包含「张」 ==========")
	zhangs := findByNameLike(db, "张")
	printUsers(zhangs)

	fmt.Println("========== 3. Where + Or 示例：age<18 或 age>=40 ==========")
	extremeAges := findYoungOrOlder(db)
	printUsers(extremeAges)

	fmt.Println("========== 4. 分页列表（page=1, size=3，create_time 倒序 + total）==========")
	page1 := listUsersByPage(db, 1, 3)
	fmt.Printf("page=%d size=%d total=%d\n", page1.Page, page1.Size, page1.Total)
	printUsers(page1.List)

	fmt.Println("========== 5. 分页列表（page=2, size=3）==========")
	page2 := listUsersByPage(db, 2, 3)
	fmt.Printf("page=%d size=%d total=%d\n", page2.Page, page2.Size, page2.Total)
	printUsers(page2.List)
}

func mustConnect() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ping 失败: %v", err)
	}
	fmt.Println("✅ GORM 已连接 go_db")
	return db
}

// prepareSchemaAndDemoData 补 gender 列，并写入演示数据（可重复运行）
func prepareSchemaAndDemoData(db *gorm.DB) {
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("迁移失败: %v", err)
	}

	demos := []User{
		{Name: "张三", Age: 22, Gender: "男", Phone: "13800001001", Password: "123456"},
		{Name: "张小花", Age: 20, Gender: "女", Phone: "13800001002", Password: "123456"},
		{Name: "李芳", Age: 25, Gender: "女", Phone: "13800001003", Password: "123456"},
		{Name: "王丽", Age: 16, Gender: "女", Phone: "13800001004", Password: "123456"}, // 女但未满 18，条件1不应命中
		{Name: "赵强", Age: 45, Gender: "男", Phone: "13800001005", Password: "123456"},
	}
	for _, u := range demos {
		var count int64
		db.Model(&User{}).Where("name = ?", u.Name).Count(&count)
		if count == 0 {
			if err := db.Create(&u).Error; err != nil {
				log.Fatalf("写入演示数据失败: %v", err)
			}
		}
	}
	// 给老数据补一个默认性别，避免空字符串干扰理解（不覆盖已有非空值）
	db.Model(&User{}).Where("gender = ? OR gender IS NULL", "").Update("gender", "男")
	fmt.Println("📦 表结构 / 演示数据已就绪")
}

// findAdultFemales 性别=女 且 age>18
// 说明：GORM 多个 Where 链式调用默认就是 AND（没有单独的 And 方法）
func findAdultFemales(db *gorm.DB) []User {
	var list []User
	err := db.Where("gender = ?", "女").
		Where("age > ?", 18). // 等价于 AND age > 18
		Order("id asc").
		Find(&list).Error
	if err != nil {
		log.Fatalf("条件查询失败: %v", err)
	}
	return list
}

// findByNameLike 姓名包含关键字（LIKE %张%）
func findByNameLike(db *gorm.DB, keyword string) []User {
	var list []User
	err := db.Where("name LIKE ?", "%"+keyword+"%").
		Order("id asc").
		Find(&list).Error
	if err != nil {
		log.Fatalf("模糊查询失败: %v", err)
	}
	return list
}

// findYoungOrOlder 演示 Or：age < 18 或 age >= 40
func findYoungOrOlder(db *gorm.DB) []User {
	var list []User
	err := db.Where("age < ?", 18).
		Or("age >= ?", 40).
		Order("age asc").
		Find(&list).Error
	if err != nil {
		log.Fatalf("Or 查询失败: %v", err)
	}
	return list
}

// listUsersByPage 分页：前端传 page、size；按 create_time 倒序，并返回 total
func listUsersByPage(db *gorm.DB, page, size int) PageResult {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	q := db.Model(&User{})

	var total int64
	if err := q.Count(&total).Error; err != nil {
		log.Fatalf("Count 失败: %v", err)
	}

	var list []User
	err := q.Order("create_time DESC").
		Limit(size).
		Offset(offset).
		Find(&list).Error
	if err != nil {
		log.Fatalf("分页查询失败: %v", err)
	}

	return PageResult{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}
}

func printUsers(list []User) {
	if len(list) == 0 {
		fmt.Println("(无数据)\n")
		return
	}
	for _, u := range list {
		fmt.Printf("  id=%-3d name=%-8s age=%-3d gender=%s create_time=%s\n",
			u.ID, u.Name, u.Age, u.Gender, u.CreateTime.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()
}

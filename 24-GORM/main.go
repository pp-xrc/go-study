/*
第 24 课：GORM 入门 —— 连接 go_db + User 模型 + CRUD

================================================================================
一、实操任务
================================================================================
  1. 安装 GORM 和 MySQL 驱动，初始化连接 go_db
  2. 定义 User 结构体，映射 user 表
  3. 用 GORM 实现：创建 / 查单个 / 查列表 / 更新 / 删除

================================================================================
二、运行
================================================================================
  go run ./24-GORM

================================================================================
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 映射 go_db.user 表（GORM 模型标签是结构体和表映射的核心）
type User struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string    `gorm:"column:name;size:64;not null"`
	Age        uint      `gorm:"column:age;not null;default:0"`
	Phone      string    `gorm:"column:phone;size:20;not null;default:''"`
	Password   string    `gorm:"column:password;size:255;not null"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
}

// TableName 指定表名（默认会推断成 users）
func (User) TableName() string {
	return "user"
}

func main() {
	db := mustConnect()

	fmt.Println("========== 1. 创建用户 ==========")
	u := createUser(db, "gorm_demo", 18, "13900000001", "123456")
	fmt.Printf("创建成功: id=%d name=%s age=%d\n\n", u.ID, u.Name, u.Age)

	fmt.Println("========== 2. 查询单个用户 ==========")
	one := getUserByID(db, u.ID)
	fmt.Printf("First 结果: %+v\n\n", one)

	fmt.Println("========== 3. 查询用户列表 ==========")
	list := listUsers(db)
	fmt.Printf("共 %d 人\n", len(list))
	for _, item := range list {
		fmt.Printf("  id=%-3d name=%-10s age=%-3d phone=%s\n",
			item.ID, item.Name, item.Age, item.Phone)
	}
	fmt.Println()

	fmt.Println("========== 4. 更新用户信息 ==========")
	updateUser(db, u.ID, 19, "13900000002")
	updated := getUserByID(db, u.ID)
	fmt.Printf("更新后: id=%d age=%d phone=%s\n\n", updated.ID, updated.Age, updated.Phone)

	fmt.Println("========== 5. 删除用户 ==========")
	deleteUser(db, u.ID)
	_, err := findUserByID(db, u.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("删除成功: id=%d 已不存在\n", u.ID)
	} else if err != nil {
		log.Fatalf("删除后查询异常: %v", err)
	} else {
		log.Fatalf("删除失败: id=%d 仍存在", u.ID)
	}
}

// mustConnect 安装驱动后，用 GORM 初始化连接 go_db
func mustConnect() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 学习阶段打印 SQL
	})
	if err != nil {
		log.Fatalf("GORM 连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层连接失败: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ping 失败: %v（请确认 MySQL 已启动、go_db 存在）", err)
	}

	fmt.Println("✅ GORM 已连接 go_db")
	return db
}

// createUser 创建用户
func createUser(db *gorm.DB, name string, age uint, phone, password string) User {
	u := User{
		Name:     name,
		Age:      age,
		Phone:    phone,
		Password: password,
	}
	if err := db.Create(&u).Error; err != nil {
		log.Fatalf("创建失败: %v", err)
	}
	return u
}

// getUserByID 查询单个用户（找不到直接 Fatal，演示用）
func getUserByID(db *gorm.DB, id uint) User {
	u, err := findUserByID(db, id)
	if err != nil {
		log.Fatalf("查询单个用户失败: %v", err)
	}
	return u
}

// findUserByID First：按主键查一条
func findUserByID(db *gorm.DB, id uint) (User, error) {
	var u User
	err := db.First(&u, id).Error
	return u, err
}

// listUsers Find：查询用户列表
func listUsers(db *gorm.DB) []User {
	var list []User
	if err := db.Order("id asc").Find(&list).Error; err != nil {
		log.Fatalf("查询列表失败: %v", err)
	}
	return list
}

// updateUser 更新用户信息
func updateUser(db *gorm.DB, id uint, age uint, phone string) {
	result := db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"age":   age,
		"phone": phone,
	})
	if result.Error != nil {
		log.Fatalf("更新失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		log.Fatalf("更新失败: id=%d 不存在", id)
	}
}

// deleteUser 删除用户
func deleteUser(db *gorm.DB, id uint) {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		log.Fatalf("删除失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		log.Fatalf("删除失败: id=%d 不存在", id)
	}
}

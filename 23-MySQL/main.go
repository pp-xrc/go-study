/*
第 23 课：Go 连接 MySQL —— 连接验证 + 查询用户列表

================================================================================
一、学习内容
================================================================================
  - database/sql 标准库 + go-sql-driver/mysql 驱动
  - DSN 连接串、Ping 验证连通性
  - Query / Scan 查询用户列表并打印

================================================================================
二、前置条件
================================================================================
  brew services start mysql
  # 库 go_db、表 user 已建好

================================================================================
三、运行
================================================================================
  go run ./23-MySQL

================================================================================
*/

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // 匿名导入：注册 mysql 驱动
)

// User 对应 go_db.user 表
type User struct {
	ID         uint
	Name       string
	Age        uint
	Phone      string
	Password   string
	CreateTime time.Time
}

func main() {
	// DSN：用户名:密码@tcp(地址)/数据库?参数
	dsn := "root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// ---------- demo2：验证连接 ----------
	if err := db.Ping(); err != nil {
		log.Fatalf("连接 MySQL 失败: %v\n请确认：1) MySQL 已启动  2) 账号密码正确  3) go_db 已创建", err)
	}

	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		log.Fatalf("查询版本失败: %v", err)
	}
	fmt.Println("✅ MySQL 连接成功")
	fmt.Printf("   驱动: github.com/go-sql-driver/mysql\n")
	fmt.Printf("   地址: 127.0.0.1:3306\n")
	fmt.Printf("   库名: go_db\n")
	fmt.Printf("   版本: %s\n\n", version)

	// 表空时插入几条演示数据，方便看查询结果
	if err := seedIfEmpty(db); err != nil {
		log.Fatalf("插入演示数据失败: %v", err)
	}

	// ---------- demo3：查询用户列表 ----------
	users, err := listUsers(db)
	if err != nil {
		log.Fatalf("查询用户列表失败: %v", err)
	}

	fmt.Printf("✅ 用户列表（共 %d 人）\n", len(users))
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("%-4s %-10s %-4s %-15s %-20s\n", "ID", "Name", "Age", "Phone", "CreateTime")
	fmt.Println("------------------------------------------------------------")
	for _, u := range users {
		fmt.Printf("%-4d %-10s %-4d %-15s %-20s\n",
			u.ID, u.Name, u.Age, u.Phone, u.CreateTime.Format("2006-01-02 15:04:05"))
	}
}

// seedIfEmpty 表为空时插入演示数据
func seedIfEmpty(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	demos := []struct {
		name, phone, password string
		age                   int
	}{
		{"张三", "13800001111", "123456", 20},
		{"李四", "13800002222", "123456", 25},
		{"王五", "13800003333", "123456", 30},
	}
	for _, d := range demos {
		_, err := db.Exec(
			"INSERT INTO user (name, age, phone, password) VALUES (?, ?, ?, ?)",
			d.name, d.age, d.phone, d.password,
		)
		if err != nil {
			return err
		}
	}
	fmt.Println("📦 已写入 3 条演示用户数据")
	return nil
}

// listUsers 用 database/sql 查询全部用户
func listUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(
		"SELECT id, name, age, phone, password, create_time FROM user ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age, &u.Phone, &u.Password, &u.CreateTime); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, rows.Err()
}

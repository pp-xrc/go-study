package model

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 本课数据库连接。controller 禁止直接使用，必须走下面的方法。
var DB *gorm.DB

type User struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string    `gorm:"column:name;size:64;not null"`
	Username   string    `gorm:"column:username;size:64;not null;default:''"`
	Nickname   string    `gorm:"column:nickname;size:64;not null;default:''"`
	Password   string    `gorm:"column:password;size:255;not null"`
	Role       string    `gorm:"column:role;size:16;not null;default:user"`
	Age        uint      `gorm:"column:age;not null;default:0"`
	Phone      string    `gorm:"column:phone;size:20;not null;default:''"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
}

func (User) TableName() string { return "user" }

func Init() error {
	dsn := "root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return err
	}
	sqlDB, err := conn.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping mysql 失败: %w", err)
	}
	if err := conn.AutoMigrate(&User{}); err != nil {
		return err
	}
	DB = conn
	log.Println("model: 已连接 go_db")
	return nil
}

// CreateUser 插入一条用户（密码应为已哈希的值）
func CreateUser(u *User) error {
	return DB.Create(u).Error
}

// GetUserByPhone 按手机号查一条；找不到返回 gorm.ErrRecordNotFound
func GetUserByPhone(phone string) (User, error) {
	var u User
	err := DB.Where("phone = ?", phone).First(&u).Error
	return u, err
}

// GetUserByID 按主键查
func GetUserByID(id uint) (User, error) {
	var u User
	err := DB.First(&u, id).Error
	return u, err
}

// GetUserList 查询有手机号的用户列表
func GetUserList() ([]User, error) {
	var list []User
	err := DB.Where("phone <> ?", "").Order("id asc").Find(&list).Error
	return list, err
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func SeedAdmin() error {
	_, err := GetUserByPhone("13800000000")
	if err == nil {
		return nil
	}
	if !IsNotFound(err) {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin := User{
		Name:     "管理员",
		Username: "admin28",
		Nickname: "管理员",
		Password: string(hashed),
		Role:     "admin",
		Phone:    "13800000000",
	}
	return CreateUser(&admin)
}

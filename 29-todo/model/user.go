package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string    `gorm:"column:name;size:64;not null"`
	Username   string    `gorm:"column:username;size:64;not null;default:''"`
	Nickname   string    `gorm:"column:nickname;size:64;not null;default:''"`
	Password   string    `gorm:"column:password;size:255;not null"`
	Role       string    `gorm:"column:role;size:16;not null;default:user"`
	Phone      string    `gorm:"column:phone;size:20;not null;default:''"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	// 一对多：一个用户有多条待办；外键是 todo.uid
	Todos []Todo `gorm:"foreignKey:UID;references:ID"`
}

func (User) TableName() string { return "user" }

func CreateUser(u *User) error {
	return DB.Create(u).Error
}

func GetUserByPhone(phone string) (User, error) {
	var u User
	err := DB.Where("phone = ?", phone).First(&u).Error
	return u, err
}

func GetUserByID(id uint) (User, error) {
	var u User
	err := DB.First(&u, id).Error
	return u, err
}

func SeedDemoUser() error {
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
	u := User{
		Name:     "演示用户",
		Username: "todo_demo",
		Nickname: "演示用户",
		Password: string(hashed),
		Role:     "user",
		Phone:    "13800000000",
	}
	return CreateUser(&u)
}

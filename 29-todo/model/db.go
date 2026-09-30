package model

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	dsn := "root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Warn),
		DisableForeignKeyConstraintWhenMigrating: true, // 旧表 user.id 是 int unsigned，避免自动建 FK 类型冲突
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
	if err := conn.AutoMigrate(&User{}, &Todo{}); err != nil {
		return err
	}
	DB = conn
	log.Println("model: 已连接 go_db，todo 表已就绪")
	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

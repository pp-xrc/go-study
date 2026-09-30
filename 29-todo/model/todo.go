package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	TodoStatusTodo = 0 // 未完成
	TodoStatusDone = 1 // 已完成
)

// Todo 对应表 todo；UID 外键关联 user.id
type Todo struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement"`
	UID        uint      `gorm:"column:uid;type:int unsigned;index;not null"`
	Title      string    `gorm:"column:title;size:128;not null"`
	Content    string    `gorm:"column:content;size:1024;not null;default:''"`
	Status     int       `gorm:"column:status;not null;default:0"` // 0 未完成 1 已完成
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
	// 多对一：这条待办属于哪个用户
	User User `gorm:"foreignKey:UID;references:ID"`
}

func (Todo) TableName() string { return "todo" }

func CreateTodo(t *Todo) error {
	return DB.Create(t).Error
}

// GetTodoListByUID 只查当前用户的待办，分页
func GetTodoListByUID(uid uint, page, size int) (list []Todo, total int64, err error) {
	q := DB.Model(&Todo{}).Where("uid = ?", uid)
	if err = q.Count(&total).Error; err != nil {
		return
	}
	err = DB.Where("uid = ?", uid).Order("id desc").Limit(size).Offset((page - 1) * size).Find(&list).Error
	return
}

// GetTodoByIDAndUID 详情：id + uid 同时匹配，别人的待办当作不存在
func GetTodoByIDAndUID(id, uid uint) (Todo, error) {
	var t Todo
	err := DB.Where("id = ? AND uid = ?", id, uid).First(&t).Error
	return t, err
}

func UpdateTodo(id, uid uint, fields map[string]any) error {
	result := DB.Model(&Todo{}).Where("id = ? AND uid = ?", id, uid).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteTodo(id, uid uint) error {
	result := DB.Where("id = ? AND uid = ?", id, uid).Delete(&Todo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

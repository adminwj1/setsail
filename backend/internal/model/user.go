package model

import (
	"time"
)

// User 用户表
type User struct {
	ID        uint      `gorm:"primaryKey;comment:用户ID" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex;not null;comment:用户名" json:"username"`
	Password  string    `gorm:"size:255;not null;comment:密码(加密存储)" json:"-"`
	Nickname  string    `gorm:"size:100;comment:昵称" json:"nickname"`
	Email     string    `gorm:"size:100;comment:邮箱" json:"email"`
	Phone     string    `gorm:"size:20;comment:手机号" json:"phone"`
	Status    int8      `gorm:"default:1;comment:状态: 1-正常, 0-禁用" json:"status"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`
}

func (User) TableName() string {
	return "sys_user"
}

// UserRole 用户角色关联表
type UserRole struct {
	ID     uint `gorm:"primaryKey;comment:主键ID" json:"id"`
	UserID uint `gorm:"not null;comment:用户ID" json:"user_id"`
	RoleID uint `gorm:"not null;comment:角色ID" json:"role_id"`
}

func (UserRole) TableName() string {
	return "sys_user_role"
}

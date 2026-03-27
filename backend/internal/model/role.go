package model

import (
	"time"
)

// Role 角色表
type Role struct {
	ID        uint      `gorm:"primaryKey;comment:角色ID" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null;comment:角色名称" json:"name"`
	Code      string    `gorm:"size:50;uniqueIndex;not null;comment:角色代码(唯一标识)" json:"code"`
	Desc      string    `gorm:"size:255;comment:角色描述" json:"desc"`
	Status    int8      `gorm:"default:1;comment:状态: 1-正常, 0-禁用" json:"status"`
	Sort      int       `gorm:"default:0;comment:排序号" json:"sort"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`
}

func (Role) TableName() string {
	return "sys_role"
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	ID     uint `gorm:"primaryKey;comment:主键ID" json:"id"`
	RoleID uint `gorm:"not null;comment:角色ID" json:"role_id"`
	MenuID uint `gorm:"not null;comment:菜单ID" json:"menu_id"`
}

func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

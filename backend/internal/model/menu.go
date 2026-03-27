package model

import (
	"time"
)

// Menu 菜单表
type Menu struct {
	ID        uint      `gorm:"primaryKey;comment:菜单ID" json:"id"`
	ParentID  uint      `gorm:"default:0;comment:父菜单ID, 0表示顶级菜单" json:"parent_id"`
	Name      string    `gorm:"size:50;not null;comment:菜单名称" json:"name"`
	Code      string    `gorm:"size:100;comment:权限代码(用于RBAC)" json:"code"`
	Path      string    `gorm:"size:200;comment:路由路径" json:"path"`
	ApiPath   string    `gorm:"size:200;comment:API路径(用于Casbin授权)" json:"api_path"`
	Icon      string    `gorm:"size:100;comment:菜单图标" json:"icon"`
	Type      int8      `gorm:"not null;comment:菜单类型: 1-目录, 2-菜单, 3-按钮" json:"type"`
	Sort      int       `gorm:"default:0;comment:排序号" json:"sort"`
	Status    int8      `gorm:"default:1;comment:状态: 1-正常, 0-禁用" json:"status"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`
	Children  []Menu    `gorm:"-" json:"children,omitempty"`
}

func (Menu) TableName() string {
	return "sys_menu"
}

const (
	MenuTypeDir    int8 = 1 // 目录
	MenuTypeMenu   int8 = 2 // 菜单
	MenuTypeButton int8 = 3 // 按钮
)

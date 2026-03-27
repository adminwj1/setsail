package model

import (
	"time"
)

// Product 产品表
type Product struct {
	ID          uint      `gorm:"primaryKey;comment:产品ID" json:"id"`
	Name        string    `gorm:"size:100;not null;comment:产品名称" json:"name"`
	Code        string    `gorm:"size:50;not null;comment:产品代号" json:"code"`
	Description string    `gorm:"type:text;comment:产品描述" json:"description"`
	Status      int8      `gorm:"default:0;comment:产品状态: 0-规划中, 1-开发中, 2-上线, 3-下线" json:"status"`
	OwnerID     uint      `gorm:"comment:负责人ID" json:"owner_id"`
	Owner       *User     `gorm:"foreignKey:OwnerID;comment:负责人" json:"owner,omitempty"`
	CreatedAt   time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"comment:更新时间" json:"updated_at"`
}

func (Product) TableName() string {
	return "pm_product"
}

const (
	ProductStatusPlanning   int8 = 0 // 规划中
	ProductStatusDeveloping int8 = 1 // 开发中
	ProductStatusOnline     int8 = 2 // 上线
	ProductStatusOffline    int8 = 3 // 下线
)

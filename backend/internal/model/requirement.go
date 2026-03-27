package model

import (
	"time"
)

// Requirement 需求表
type Requirement struct {
	ID          uint        `gorm:"primaryKey;comment:需求ID" json:"id"`
	ProductID   uint        `gorm:"not null;comment:所属产品ID" json:"product_id"`
	Product     *Product    `gorm:"foreignKey:ProductID;comment:所属产品" json:"product,omitempty"`
	Title       string      `gorm:"size:200;not null;comment:需求标题" json:"title"`
	Type        int8        `gorm:"not null;comment:需求类型: 1-业务需求, 2-用户需求, 3-研发需求" json:"type"`
	Priority    int8        `gorm:"default:2;comment:优先级: 1-高, 2-中, 3-低" json:"priority"`
	Status      int8        `gorm:"default:0;comment:需求状态: 0-待评审, 1-已采纳, 2-开发中, 3-已完成, 4-已拒绝" json:"status"`
	Description string      `gorm:"type:text;comment:需求描述" json:"description"`
	CreatorID   uint        `gorm:"comment:创建人ID" json:"creator_id"`
	Creator     *User       `gorm:"foreignKey:CreatorID;comment:创建人" json:"creator,omitempty"`
	AssigneeID  uint        `gorm:"comment:负责人ID" json:"assignee_id"`
	Assignee    *User       `gorm:"foreignKey:AssigneeID;comment:负责人" json:"assignee,omitempty"`
	CreatedAt   time.Time   `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"comment:更新时间" json:"updated_at"`
}

func (Requirement) TableName() string {
	return "pm_requirement"
}

const (
	ReqTypeBusiness    int8 = 1 // 业务需求
	ReqTypeUser        int8 = 2 // 用户需求
	ReqTypeDevelopment int8 = 3 // 研发需求
)

const (
	ReqPriorityHigh   int8 = 1 // 高
	ReqPriorityMedium int8 = 2 // 中
	ReqPriorityLow    int8 = 3 // 低
)

const (
	ReqStatusPending    int8 = 0 // 待评审
	ReqStatusAccepted   int8 = 1 // 已采纳
	ReqStatusDeveloping int8 = 2 // 开发中
	ReqStatusCompleted  int8 = 3 // 已完成
	ReqStatusRejected   int8 = 4 // 已拒绝
)

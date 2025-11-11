// Package model 定义了与数据库表对应的 GORM 模型
package model

import (
	"time"

	"gorm.io/datatypes"
)

// Submission 对应于数据库中的 `submissions` 表
type Submission struct {
	ID              uint           `gorm:"primarykey"`
	FormID          uint           `gorm:"not null"`
	SubmitterID     *uint          `gorm:"null"`     // 提交者ID, 允许匿名
	Data            datatypes.JSON `gorm:"not null"` // 用户提交的答案数据
	RawScore        *int           `gorm:"null"`     // 原始得分
	MaxScore        *int           `gorm:"null"`     // 总分
	DurationSeconds *uint          `gorm:"null"`     // 答题用时
	ClientIP        string         `gorm:"type:varchar(45)"`
	UserAgent       string         `gorm:"type:text"`
	CreatedAt       time.Time

	// 定义关联关系
	Form      Form `gorm:"foreignKey:FormID"`
	Submitter User `gorm:"foreignKey:SubmitterID"`
}

// TableName 指定 Submission 模型对应的数据库表名
func (Submission) TableName() string {
	return "submissions"
}

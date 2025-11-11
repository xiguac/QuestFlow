// Package model 定义了与数据库表对应的 GORM 模型
package model

import (
	"time"

	"github.com/teris-io/shortid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Form 对应于数据库中的 `forms` 表
type Form struct {
	ID          uint           `gorm:"primarykey"`
	FormKey     string         `gorm:"type:varchar(20);uniqueIndex;not null"`
	CreatorID   uint           `gorm:"not null"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Definition  datatypes.JSON `gorm:"not null"` // 使用 GORM 的 JSON 类型
	Status      uint8          `gorm:"type:tinyint unsigned;not null;default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// 定义关联关系
	Creator User `gorm:"foreignKey:CreatorID"`
}

// TableName 指定 Form 模型对应的数据库表名
func (Form) TableName() string {
	return "forms"
}

// BeforeCreate 是一个 GORM Hook, 会在创建记录到数据库之前被调用
func (f *Form) BeforeCreate(tx *gorm.DB) (err error) {
	// 如果 FormKey 为空，则生成一个新的
	if f.FormKey == "" {
		sid, err := shortid.New(1, shortid.DefaultABC, 2342)
		if err != nil {
			return err
		}
		f.FormKey, err = sid.Generate()
		if err != nil {
			return err
		}
	}
	return
}

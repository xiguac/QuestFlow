// Package model 定义了与数据库表对应的 GORM 模型
package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 对应于数据库中的 `users` 表
type User struct {
	ID           uint           `gorm:"primarykey"`
	Username     string         `gorm:"type:varchar(50);uniqueIndex;not null"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	Email        *string        `gorm:"type:varchar(100);uniqueIndex"`
	Nickname     string         `gorm:"type:varchar(50);default:''"`
	Avatar       string         `gorm:"type:varchar(255);default:''"`
	Role         uint8          `gorm:"type:tinyint unsigned;not null;default:1"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName 指定 User 模型对应的数据库表名
func (User) TableName() string {
	return "users"
}

// SetPassword 对密码进行 bcrypt 哈希加密
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(bytes)
	return nil
}

// CheckPassword 验证密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

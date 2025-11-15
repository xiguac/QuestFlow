// Package repository 封装了数据访问逻辑
package repository

import (
	"questflow/internal/model"

	"gorm.io/gorm"
)

// UserRepository 定义了用户数据仓库的接口
type UserRepository interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Update(user *model.User) error
}

// userGormRepository 是 UserRepository 的 GORM 实现
type userGormRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建一个新的 UserRepository 实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userGormRepository{db: db}
}

// Create 在数据库中创建一个新用户
func (r *userGormRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByUsername 通过用户名查找用户
func (r *userGormRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 通过用户 ID 查找用户
func (r *userGormRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *userGormRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Package repository 封装了数据访问逻辑
package repository

import (
	"questflow/internal/model"

	"gorm.io/gorm"
)

// FormRepository 定义了表单数据仓库的接口
type FormRepository interface {
	Create(form *model.Form) error
	FindByKey(key string) (*model.Form, error)
	FindByID(id uint) (*model.Form, error)
	FindByCreatorID(creatorID uint) ([]model.Form, error)
	Delete(form *model.Form) error
	Update(form *model.Form) error
}

// formGormRepository 是 FormRepository 的 GORM 实现
type formGormRepository struct {
	db *gorm.DB
}

// NewFormRepository 创建一个新的 FormRepository 实例
func NewFormRepository(db *gorm.DB) FormRepository {
	return &formGormRepository{db: db}
}

// Create 在数据库中创建一个新表单
func (r *formGormRepository) Create(form *model.Form) error {
	return r.db.Create(form).Error
}

// FindByKey 通过唯一的 form_key 查找表单
func (r *formGormRepository) FindByKey(key string) (*model.Form, error) {
	var form model.Form
	err := r.db.Where("form_key = ?", key).First(&form).Error
	if err != nil {
		return nil, err
	}
	return &form, nil
}

// FindByID 通过主键 ID 查找表单
func (r *formGormRepository) FindByID(id uint) (*model.Form, error) {
	var form model.Form
	// Preload("Creator") 可以在查询表单时，同时加载关联的创建者信息
	err := r.db.Preload("Creator").First(&form, id).Error
	if err != nil {
		return nil, err
	}
	return &form, nil
}

// FindByCreatorID 查找某个用户创建的所有表单
func (r *formGormRepository) FindByCreatorID(creatorID uint) ([]model.Form, error) {
	var forms []model.Form
	err := r.db.Where("creator_id = ?", creatorID).Order("created_at desc").Find(&forms).Error
	return forms, err
}

// Delete 从数据库中删除一个表单 (GORM 会自动处理软删除)
func (r *formGormRepository) Delete(form *model.Form) error {
	return r.db.Delete(form).Error
}

// Update 更新数据库中的表单信息
func (r *formGormRepository) Update(form *model.Form) error {
	return r.db.Save(form).Error
}

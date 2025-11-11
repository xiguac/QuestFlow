// Package repository 封装了数据访问逻辑
package repository

import (
	"questflow/internal/model"

	"gorm.io/gorm"
)

// SubmissionRepository 定义了提交数据仓库的接口
type SubmissionRepository interface {
	Create(submission *model.Submission) error
	FindByFormID(formID uint) ([]model.Submission, error)
}

// submissionGormRepository 是 SubmissionRepository 的 GORM 实现
type submissionGormRepository struct {
	db *gorm.DB
}

// NewSubmissionRepository 创建一个新的 SubmissionRepository 实例
func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionGormRepository{db: db}
}

// Create 在数据库中创建一条新的提交记录
func (r *submissionGormRepository) Create(submission *model.Submission) error {
	return r.db.Create(submission).Error
}

// FindByFormID 查找某个表单下的所有提交记录
func (r *submissionGormRepository) FindByFormID(formID uint) ([]model.Submission, error) {
	var submissions []model.Submission
	err := r.db.Where("form_id = ?", formID).Order("created_at desc").Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

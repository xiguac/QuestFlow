// Package repository 封装了数据访问逻辑
package repository

import (
	"fmt"
	"questflow/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

// FilterCondition 定义了单个筛选条件
type FilterCondition struct {
	QuestionID   string   `json:"questionId"`
	QuestionType string   `json:"questionType"`
	Operator     string   `json:"operator"` // "equals", "not_equals", "contains", "not_contains"
	Value        []string `json:"value"`    // 答案值，使用数组以支持多选
}

// SubmissionRepository 接口定义
type SubmissionRepository interface {
	Create(submission *model.Submission) error
	FindByFormID(formID uint) ([]model.Submission, error)
	FindWithFilters(formID uint, startTime, endTime *time.Time, conditions []FilterCondition) ([]model.Submission, error)
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
	err := r.db.Where("form_id = ?", formID).Order("created_at asc").Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

// FindWithFilters 【核心修复】使用更健壮的 SQL 构建逻辑
func (r *submissionGormRepository) FindWithFilters(formID uint, startTime, endTime *time.Time, conditions []FilterCondition) ([]model.Submission, error) {
	var submissions []model.Submission

	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("SELECT * FROM submissions WHERE form_id = ?")

	args := []interface{}{formID}

	if startTime != nil {
		sqlBuilder.WriteString(" AND created_at >= ?")
		args = append(args, *startTime)
	}
	if endTime != nil {
		sqlBuilder.WriteString(" AND created_at <= ?")
		args = append(args, *endTime)
	}

	for _, cond := range conditions {
		jsonPath := fmt.Sprintf("$.\"%s\"", cond.QuestionID)

		switch cond.QuestionType {
		case "single_choice", "judgment", "text_input":
			// 对于这些题型，答案是一个字符串
			value := cond.Value[0] // 单值
			switch cond.Operator {
			case "equals":
				// JSON_EXTRACT 返回带引号的字符串, JSON_UNQUOTE 去掉引号
				sqlBuilder.WriteString(fmt.Sprintf(" AND JSON_UNQUOTE(data->'%s') = ?", jsonPath))
				args = append(args, value)
			case "not_equals":
				// 检查值不等于或者该键不存在
				sqlBuilder.WriteString(fmt.Sprintf(" AND (data->'%s' IS NULL OR JSON_UNQUOTE(data->'%s') != ?)", jsonPath, jsonPath))
				args = append(args, value)
			}
		case "multi_choice":
			// 对于多选题，答案是一个数组
			// 构造一个 JSON 数组字符串用于查询, e.g., '["val1", "val2"]'
			var valuePlaceholders []string
			for range cond.Value {
				valuePlaceholders = append(valuePlaceholders, "?")
			}
			jsonArrayForQuery := fmt.Sprintf("CAST(JSON_ARRAY(%s) AS JSON)", strings.Join(valuePlaceholders, ","))

			for _, v := range cond.Value {
				args = append(args, v)
			}

			switch cond.Operator {
			case "contains": // 包含 cond.Value 中的任意一个
				sqlBuilder.WriteString(fmt.Sprintf(" AND JSON_OVERLAPS(data->'%s', %s)", jsonPath, jsonArrayForQuery))
			case "not_contains": // 不包含 cond.Value 中的任何一个
				sqlBuilder.WriteString(fmt.Sprintf(" AND NOT JSON_OVERLAPS(data->'%s', %s)", jsonPath, jsonArrayForQuery))
			case "equals": // 完全匹配（忽略顺序）
				// 检查两个数组长度是否相等 并且 数据库中的数组完全包含查询数组的所有元素
				sqlBuilder.WriteString(fmt.Sprintf(" AND JSON_LENGTH(data->'%s') = ? AND JSON_CONTAINS(data->'%s', %s)", jsonPath, jsonPath, jsonArrayForQuery))
				args = append(args, len(cond.Value)) // 添加数组长度作为参数
			}
		}
	}

	sqlBuilder.WriteString(" ORDER BY created_at asc")

	err := r.db.Raw(sqlBuilder.String(), args...).Scan(&submissions).Error
	if err != nil {
		return nil, err
	}

	return submissions, nil
}

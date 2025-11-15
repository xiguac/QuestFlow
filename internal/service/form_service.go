// Package service 包含了应用的业务逻辑
package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"questflow/internal/model"
	"questflow/internal/repository"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// --- 为统计结果定义清晰的结构体 ---

// OptionStat 存储单个选项的统计
type OptionStat struct {
	Text  string `json:"text"`
	Count int    `json:"count"`
}

// QuestionStat 存储单个问题的统计结果
type QuestionStat struct {
	QuestionID   string       `json:"question_id"`
	QuestionType string       `json:"question_type"`
	Title        string       `json:"title"`
	OptionStats  []OptionStat `json:"option_stats,omitempty"` // 用于选择题
	TextAnswers  []string     `json:"text_answers,omitempty"` // 用于填空题
}

// FormStats 最终返回给前端的完整统计数据结构
type FormStats struct {
	TotalSubmissions int            `json:"total_submissions"`
	QuestionStats    []QuestionStat `json:"question_stats"`
}

// --- 表单定义解析用的临时结构体 ---
type questionDefinition struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Options []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"options"`
}

type formDefinition struct {
	Questions []questionDefinition `json:"questions"`
}

// --- 更新 Service 接口和实现 ---
type FormService interface {
	CreateForm(creatorID uint, title, description string, definition datatypes.JSON) (*model.Form, error)
	GetPublicFormByKey(key string) (*model.Form, error)
	GetFormStatistics(formID uint, userID uint) (*FormStats, error)
	GetFormsByCreator(userID uint) ([]model.Form, error)
	DeleteForm(formID uint, userID uint) error
	GetFormForEditing(formID uint, userID uint) (*model.Form, error)
	UpdateForm(formID, userID uint, title, description string, definition datatypes.JSON) (*model.Form, error)
	UpdateFormStatus(formID, userID uint, status uint8) error
	ExportFormSubmissions(formID, userID uint, startTime, endTime *time.Time, conditions []repository.FilterCondition) (*bytes.Buffer, *model.Form, error)
}

type formServiceImpl struct {
	formRepo       repository.FormRepository
	submissionRepo repository.SubmissionRepository
	excelService   ExcelService
}

func NewFormService(formRepo repository.FormRepository, submissionRepo repository.SubmissionRepository) FormService {
	return &formServiceImpl{
		formRepo:       formRepo,
		submissionRepo: submissionRepo,
		excelService:   NewExcelService(),
	}
}

// CreateForm
func (s *formServiceImpl) CreateForm(creatorID uint, title, description string, definition datatypes.JSON) (*model.Form, error) {
	newForm := &model.Form{
		CreatorID:   creatorID,
		Title:       title,
		Description: description,
		Definition:  definition,
		Status:      1, // 1: 草稿
	}
	err := s.formRepo.Create(newForm)
	if err != nil {
		return nil, err
	}
	return newForm, nil
}

// UpdateForm
func (s *formServiceImpl) UpdateForm(formID, userID uint, title, description string, definition datatypes.JSON) (*model.Form, error) {
	form, err := s.GetFormForEditing(formID, userID) // 复用权限检查逻辑
	if err != nil {
		return nil, err
	}

	form.Title = title
	form.Description = description
	form.Definition = definition

	if err := s.formRepo.Update(form); err != nil {
		return nil, err
	}
	return form, nil
}

// UpdateFormStatus
func (s *formServiceImpl) UpdateFormStatus(formID, userID uint, status uint8) error {
	form, err := s.GetFormForEditing(formID, userID) // 复用权限检查逻辑
	if err != nil {
		return err
	}
	// 简单校验一下状态值
	if status < 1 || status > 3 {
		return errors.New("invalid status value")
	}

	form.Status = status
	return s.formRepo.Update(form)
}

// GetFormForEditing
func (s *formServiceImpl) GetFormForEditing(formID uint, userID uint) (*model.Form, error) {
	form, err := s.formRepo.FindByID(formID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("form not found")
		}
		return nil, err
	}
	if form.CreatorID != userID {
		return nil, errors.New("access denied")
	}
	return form, nil
}

// GetPublicFormByKey
func (s *formServiceImpl) GetPublicFormByKey(key string) (*model.Form, error) {
	form, err := s.formRepo.FindByKey(key)
	if err != nil {
		return nil, err
	}
	// 只有已发布的问卷才能公开访问
	if form.Status != 2 {
		return nil, errors.New("form not available")
	}
	return form, nil
}

// GetFormsByCreator 获取用户创建的表单列表
func (s *formServiceImpl) GetFormsByCreator(userID uint) ([]model.Form, error) {
	return s.formRepo.FindByCreatorID(userID)
}

// DeleteForm 处理删除表单的业务逻辑
func (s *formServiceImpl) DeleteForm(formID uint, userID uint) error {
	form, err := s.formRepo.FindByID(formID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("form not found")
		}
		return err
	}
	if form.CreatorID != userID {
		return errors.New("access denied")
	}
	return s.formRepo.Delete(form)
}

// ExportFormSubmissions 实现导出业务逻辑
func (s *formServiceImpl) ExportFormSubmissions(formID, userID uint, startTime, endTime *time.Time, conditions []repository.FilterCondition) (*bytes.Buffer, *model.Form, error) {
	// 1. 获取表单并进行权限验证
	form, err := s.formRepo.FindByID(formID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("form not found")
		}
		return nil, nil, err
	}
	if form.CreatorID != userID {
		return nil, nil, errors.New("access denied")
	}

	// 2. 解析表单定义
	var def formDefinition
	if err := json.Unmarshal(form.Definition, &def); err != nil {
		return nil, nil, errors.New("failed to parse form definition")
	}

	// 3. 根据筛选条件获取提交记录
	submissions, err := s.submissionRepo.FindWithFilters(formID, startTime, endTime, conditions)
	if err != nil {
		return nil, nil, err
	}

	if len(submissions) == 0 {
		return nil, nil, errors.New("no submissions found for the given criteria")
	}

	// 4. 调用 ExcelService 生成文件流
	buffer, err := s.excelService.ExportSubmissionsToExcel(&def, submissions)
	if err != nil {
		return nil, nil, err
	}

	return buffer, form, nil
}

// GetFormStatistics
func (s *formServiceImpl) GetFormStatistics(formID uint, userID uint) (*FormStats, error) {
	// 1. 获取表单并进行权限验证
	form, err := s.formRepo.FindByID(formID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("form not found")
		}
		return nil, err
	}
	if form.CreatorID != userID {
		return nil, errors.New("access denied")
	}

	// 2. 获取该表单的所有提交记录
	submissions, err := s.submissionRepo.FindByFormID(formID)
	if err != nil {
		return nil, err
	}

	// 3. 解析表单定义
	var def formDefinition
	if err := json.Unmarshal(form.Definition, &def); err != nil {
		return nil, errors.New("failed to parse form definition")
	}

	questionDefMap := make(map[string]questionDefinition)
	optionTextMap := make(map[string]string)
	for _, q := range def.Questions {
		questionDefMap[q.ID] = q
		for _, opt := range q.Options {
			optionTextMap[opt.ID] = opt.Text
		}
	}

	// 4. 初始化统计结果结构
	statsResult := &FormStats{
		TotalSubmissions: len(submissions),
		QuestionStats:    make([]QuestionStat, 0, len(def.Questions)),
	}
	tempStats := make(map[string]map[string]int) // {questionID: {optionID: count}}
	tempTextAnswers := make(map[string][]string) // {questionID: [answer1, answer2]}

	// 5. 遍历所有提交记录，进行数据聚合
	for _, sub := range submissions {
		var answers map[string]interface{}
		if err := json.Unmarshal(sub.Data, &answers); err != nil {
			continue
		}

		for qID, ans := range answers {
			qDef, ok := questionDefMap[qID]
			if !ok {
				continue
			}

			// 确保为每个问题初始化统计map
			if _, exists := tempStats[qID]; !exists {
				tempStats[qID] = make(map[string]int)
			}

			switch qDef.Type {
			// 单选题和判断题的答案是 string
			case "single_choice", "judgment":
				if optID, ok := ans.(string); ok {
					tempStats[qID][optID]++
				}
			// 多选题的答案是 []string
			case "multi_choice":
				if opts, ok := ans.([]interface{}); ok {
					for _, opt := range opts {
						if optID, ok := opt.(string); ok {
							tempStats[qID][optID]++
						}
					}
				}
			// 填空题的答案是 string
			case "text_input":
				if text, ok := ans.(string); ok {
					tempTextAnswers[qID] = append(tempTextAnswers[qID], text)
				}
			}
		}
	}

	// 6. 将聚合后的数据整理成最终的返回格式
	for _, qDef := range def.Questions {
		qStat := QuestionStat{
			QuestionID:   qDef.ID,
			QuestionType: qDef.Type,
			Title:        qDef.Title,
		}

		// 统一处理所有基于选项的题型
		if qDef.Type == "single_choice" || qDef.Type == "multi_choice" || qDef.Type == "judgment" {
			// 初始化所有选项的计数为0
			optionCounts := make(map[string]int)
			for _, opt := range qDef.Options {
				optionCounts[opt.ID] = 0
			}
			// 如果有统计数据，则填充
			if counts, ok := tempStats[qDef.ID]; ok {
				for optID, count := range counts {
					optionCounts[optID] = count
				}
			}
			// 转换为最终的数组格式
			qStat.OptionStats = make([]OptionStat, 0, len(qDef.Options))
			for _, opt := range qDef.Options {
				qStat.OptionStats = append(qStat.OptionStats, OptionStat{
					Text:  optionTextMap[opt.ID],
					Count: optionCounts[opt.ID],
				})
			}
		} else if qDef.Type == "text_input" {
			qStat.TextAnswers = tempTextAnswers[qDef.ID]
		}

		statsResult.QuestionStats = append(statsResult.QuestionStats, qStat)
	}

	return statsResult, nil
}

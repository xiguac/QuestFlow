// Package service 包含了应用的业务逻辑
package service

import (
	"context"
	"encoding/json"
	"errors"
	"questflow/internal/model"
	"questflow/internal/repository"
	"questflow/pkg/redis" // 只导入我们自己的 redis 包
	"time"

	"gorm.io/datatypes"
)

// SubmissionMessage 定义了发送到消息队列的提交数据的结构
type SubmissionMessage struct {
	FormID      uint            `json:"form_id"`
	Data        json.RawMessage `json:"data"`
	ClientIP    string          `json:"client_ip"`
	UserAgent   string          `json:"user_agent"`
	SubmitterID *uint           `json:"submitter_id,omitempty"`
	SubmittedAt time.Time       `json:"submitted_at"`
}

// SubmissionService 定义了提交服务的接口
type SubmissionService interface {
	CreateSubmission(form *model.Form, data datatypes.JSON, clientIP string, userAgent string, submitterID *uint) (string, error)
	ProcessSubmission(msg SubmissionMessage) error
}

// submissionServiceImpl 是 SubmissionService 的实现
type submissionServiceImpl struct {
	submissionRepo repository.SubmissionRepository
}

// NewSubmissionService 创建一个新的 SubmissionService 实例
func NewSubmissionService(repo repository.SubmissionRepository) SubmissionService {
	return &submissionServiceImpl{submissionRepo: repo}
}

// CreateSubmission (生产者逻辑): 调用封装好的 Redis 发布方法
func (s *submissionServiceImpl) CreateSubmission(form *model.Form, data datatypes.JSON, clientIP string, userAgent string, submitterID *uint) (string, error) {
	// 1. 业务校验 (在 Web 服务中快速完成)
	if form.Status != 2 { // 假设 2 代表 "已发布"
		return "", errors.New("form is not published")
	}

	// 2. 构造消息
	msg := SubmissionMessage{
		FormID:      form.ID,
		Data:        json.RawMessage(data),
		ClientIP:    clientIP,
		UserAgent:   userAgent,
		SubmitterID: submitterID,
		SubmittedAt: time.Now(),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return "", errors.New("failed to serialize submission message")
	}

	// 3. 调用我们封装好的方法发送消息
	messageID, err := redis.PublishSubmissionMessage(context.Background(), msgBytes)
	if err != nil {
		return "", errors.New("failed to publish submission message to stream")
	}

	return messageID, nil
}

// ProcessSubmission (消费者逻辑): 包含了写入数据库的逻辑
func (s *submissionServiceImpl) ProcessSubmission(msg SubmissionMessage) error {
	if s.submissionRepo == nil {
		return errors.New("submission repository is not initialized")
	}

	newSubmission := &model.Submission{
		FormID:      msg.FormID,
		SubmitterID: msg.SubmitterID,
		Data:        datatypes.JSON(msg.Data),
		ClientIP:    msg.ClientIP,
		UserAgent:   msg.UserAgent,
		CreatedAt:   msg.SubmittedAt,
	}

	return s.submissionRepo.Create(newSubmission)
}

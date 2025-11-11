// Package handler 存放 HTTP 请求的处理器函数
package handler

import (
	"encoding/json"
	"net/http"
	"questflow/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type SubmissionHandler struct {
	submissionService service.SubmissionService
	formService       service.FormService
}

func NewSubmissionHandler(subService service.SubmissionService, formService service.FormService) *SubmissionHandler {
	return &SubmissionHandler{
		submissionService: subService,
		formService:       formService,
	}
}

type CreateSubmissionRequest struct {
	Data json.RawMessage `json:"data" binding:"required"`
}

// CreateSubmission 处理提交表单数据的请求
func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	formKey := c.Param("form_key")
	if formKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "缺少 form_key"})
		return
	}

	var req CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": err.Error()})
		return
	}

	form, err := h.formService.GetPublicFormByKey(formKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "表单不存在"})
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	var submitterID *uint

	// 调用改造后的 service 方法
	messageID, err := h.submissionService.CreateSubmission(form, datatypes.JSON(req.Data), clientIP, userAgent, submitterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5000, "message": "提交失败", "error": err.Error()})
		return
	}

	// 快速返回成功响应，现在返回的是消息ID
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "提交成功，正在处理中...",
		"data": gin.H{
			"message_id": messageID,
		},
	})
}

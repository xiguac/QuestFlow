// Package handler 存放 HTTP 请求的处理器函数
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"questflow/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FormHandler 封装了表单相关的 HTTP 处理器
type FormHandler struct {
	formService service.FormService
}

// NewFormHandler 创建一个新的 FormHandler
func NewFormHandler(formService service.FormService) *FormHandler {
	return &FormHandler{formService: formService}
}

// CreateFormRequest 定义了创建/更新表单请求的 JSON 结构体
type UpsertFormRequest struct {
	Title       string          `json:"title" binding:"required,max=255"`
	Description string          `json:"description"`
	Definition  json.RawMessage `json:"definition" binding:"required"`
}

// CreateForm 处理创建新表单的请求
func (h *FormHandler) CreateForm(c *gin.Context) {
	var req UpsertFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": err.Error()})
		return
	}

	claims, _ := c.Get("user_claims")
	userClaims := claims.(*service.CustomClaims)
	creatorID := userClaims.UserID

	definitionJSON := datatypes.JSON(req.Definition)

	form, err := h.formService.CreateForm(creatorID, req.Title, req.Description, definitionJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5000, "message": "创建表单失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    form,
	})
}

// UpdateForm 处理更新表单的请求
func (h *FormHandler) UpdateForm(c *gin.Context) {
	formID, err := getFormIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "无效的 form_id"})
		return
	}

	var req UpsertFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": err.Error()})
		return
	}

	claims, _ := c.Get("user_claims")
	userClaims := claims.(*service.CustomClaims)

	definitionJSON := datatypes.JSON(req.Definition)

	form, err := h.formService.UpdateForm(formID, userClaims.UserID, req.Title, req.Description, definitionJSON)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    form,
	})
}

// UpdateFormStatusRequest 定义了更新状态请求的 body
type UpdateFormStatusRequest struct {
	Status uint8 `json:"status" binding:"required,min=1,max=3"`
}

// UpdateFormStatus 处理更新表单状态的请求
func (h *FormHandler) UpdateFormStatus(c *gin.Context) {
	formID, err := getFormIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "无效的 form_id"})
		return
	}

	var req UpdateFormStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": err.Error()})
		return
	}

	claims, _ := c.Get("user_claims")
	userClaims := claims.(*service.CustomClaims)

	err = h.formService.UpdateFormStatus(formID, userClaims.UserID, req.Status)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "状态更新成功",
	})
}

// GetFormDetails 处理获取单个表单详细信息以供编辑的请求
func (h *FormHandler) GetFormDetails(c *gin.Context) {
	formID, err := getFormIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "无效的 form_id"})
		return
	}
	claims, _ := c.Get("user_claims")
	userClaims := claims.(*service.CustomClaims)

	form, err := h.formService.GetFormForEditing(formID, userClaims.UserID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    form,
	})
}

// GetPublicForm 处理获取公开表单定义的请求
func (h *FormHandler) GetPublicForm(c *gin.Context) {
	formKey := c.Param("form_key")
	if formKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "缺少 form_key"})
		return
	}

	form, err := h.formService.GetPublicFormByKey(formKey)
	if err != nil {
		if err.Error() == "form not available" {
			c.JSON(http.StatusForbidden, gin.H{"code": 4003, "message": "该问卷未发布或已关闭"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "表单未找到"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"form_key":    form.FormKey,
			"title":       form.Title,
			"description": form.Description,
			"definition":  form.Definition,
		},
	})
}

// GetStatistics 处理获取表单统计数据的请求
func (h *FormHandler) GetStatistics(c *gin.Context) {
	formID, err := getFormIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "无效的 form_id"})
		return
	}

	claims, _ := c.Get("user_claims")
	userClaims := claims.(*service.CustomClaims)
	userID := userClaims.UserID

	stats, err := h.formService.GetFormStatistics(formID, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// GetMyForms 处理获取当前用户创建的表单列表的请求
func (h *FormHandler) GetMyForms(c *gin.Context) {
	claims, _ := c.Get("user_claims")
	userClaims := claims.(*service.CustomClaims)
	userID := userClaims.UserID

	forms, err := h.formService.GetFormsByCreator(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5000, "message": "获取表单列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    forms,
	})
}

// DeleteForm 处理删除表单的请求
func (h *FormHandler) DeleteForm(c *gin.Context) {
	formID, err := getFormIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "无效的 form_id"})
		return
	}

	claims, _ := c.Get("user_claims")
	userClaims := claims.(*service.CustomClaims)
	userID := userClaims.UserID

	err = h.formService.DeleteForm(formID, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// --- 辅助函数 ---

func getFormIDFromParam(c *gin.Context) (uint, error) {
	formIDStr := c.Param("form_id")
	formID, err := strconv.ParseUint(formIDStr, 10, 32)
	return uint(formID), err
}

func handleServiceError(c *gin.Context, err error) {
	switch err.Error() {
	case "access denied":
		c.JSON(http.StatusForbidden, gin.H{"code": 4003, "message": "无权操作此表单"})
	case "form not found":
		c.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "表单未找到"})
	case "invalid status value":
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "无效的状态值"})
	default:
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "记录未找到"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5000, "message": "服务器内部错误", "error": err.Error()})
	}
}

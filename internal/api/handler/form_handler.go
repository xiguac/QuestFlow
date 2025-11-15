// Package handler 存放 HTTP 请求的处理器函数
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"questflow/internal/repository" // 引入 repository 以使用 FilterCondition
	"questflow/internal/service"
	"strconv"
	"time"

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

// UpsertFormRequest 定义了创建/更新表单请求的 JSON 结构体
type UpsertFormRequest struct {
	Title       string          `json:"title" binding:"required,max=255"`
	Description string          `json:"description"`
	Definition  json.RawMessage `json:"definition" binding:"required"`
}

// UpdateFormStatusRequest 定义了更新状态请求的 body
type UpdateFormStatusRequest struct {
	Status uint8 `json:"status" binding:"required,min=1,max=3"`
}

// ExportRequest 定义了导出请求的 JSON 结构体
type ExportRequest struct {
	StartTime  *time.Time                   `json:"startTime"`
	EndTime    *time.Time                   `json:"endTime"`
	Conditions []repository.FilterCondition `json:"conditions"`
}

// ExportSubmissions 处理导出提交数据的请求 (重构为 POST)
func (h *FormHandler) ExportSubmissions(c *gin.Context) {
	formID, err := getFormIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "无效的 form_id"})
		return
	}

	var req ExportRequest
	// 绑定请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "无效的筛选条件格式: " + err.Error()})
		return
	}

	claims, _ := c.Get("user_claims")
	userClaims := claims.(*service.CustomClaims)

	// 调用 service 获取 Excel 文件流
	buffer, form, err := h.formService.ExportFormSubmissions(formID, userClaims.UserID, req.StartTime, req.EndTime, req.Conditions)
	if err != nil {
		if err.Error() == "no submissions found for the given criteria" {
			c.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "在指定条件下未找到任何提交数据"})
			return
		}
		handleServiceError(c, err)
		return
	}

	// 设置 HTTP 响应头，告知浏览器这是一个文件下载
	fileName := fmt.Sprintf("%s-submissions-%s.xlsx", form.Title, time.Now().Format("200601021504"))
	encodedFileName := url.QueryEscape(fileName)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+encodedFileName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", fmt.Sprintf("%d", buffer.Len()))

	// 将文件流写入响应体
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())
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
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": form})
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
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": form})
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
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "状态更新成功"})
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
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": form})
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
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"form_key": form.FormKey, "title": form.Title, "description": form.Description, "definition": form.Definition}})
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
	stats, err := h.formService.GetFormStatistics(formID, userClaims.UserID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": stats})
}

// GetMyForms 处理获取当前用户创建的表单列表的请求
func (h *FormHandler) GetMyForms(c *gin.Context) {
	claims, _ := c.Get("user_claims")
	userClaims := claims.(*service.CustomClaims)
	forms, err := h.formService.GetFormsByCreator(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5000, "message": "获取表单列表失败", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": forms})
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
	err = h.formService.DeleteForm(formID, userClaims.UserID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
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

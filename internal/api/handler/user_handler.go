// Package handler 存放 HTTP 请求的处理器函数
package handler

import (
	"net/http"
	"questflow/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 封装了用户相关的 HTTP 处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建一个新的 UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterRequest 定义了注册请求的 JSON 结构体
type RegisterRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=20"`
	Password string  `json:"password" binding:"required,min=6,max=30"`
	Email    *string `json:"email" binding:"omitempty,email"`
}

// Register 处理用户注册请求
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": err.Error(), "data": nil})
		return
	}
	user, err := h.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// LoginRequest 定义了登录请求的 JSON 结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"binding:"required"`
}

// Login 处理用户登录请求
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": err.Error(), "data": nil})
		return
	}
	token, user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 4001, "message": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": gin.H{
			"token": token,
			"user_info": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"nickname": user.Nickname,
				"role":     user.Role,
			},
		},
	})
}

// UpdateProfileRequest 定义了更新用户信息请求的 JSON 结构体
type UpdateProfileRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=20"`
	Password string `json:"password" binding:"omitempty,min=6,max=30"`
}

// UpdateProfile 处理更新用户信息的请求
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": err.Error()})
		return
	}

	// 确保至少提供了一项要更新的内容
	if req.Username == "" && req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": "no fields to update"})
		return
	}

	claims, _ := c.Get("user_claims")
	userClaims := claims.(*service.CustomClaims)

	_, err := h.userService.UpdateProfile(userClaims.UserID, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5000, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "用户信息更新成功，请重新登录以使新信息生效。",
	})
}

// GetCurrentUser 获取当前登录用户的信息
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	claims, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "无法获取用户信息",
		})
		return
	}
	userClaims, ok := claims.(*service.CustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "用户信息格式错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":       userClaims.UserID,
			"username": userClaims.Username,
			"role":     userClaims.Role,
		},
	})
}

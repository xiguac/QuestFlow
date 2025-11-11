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
	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": err.Error(), "data": nil})
		return
	}

	// 2. 调用 service 层处理业务逻辑
	user, err := h.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		// 根据错误类型返回不同的状态码和信息
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error(), "data": nil})
		return
	}

	// 3. 构造并返回成功响应
	// 注意：实际项目中不应返回密码哈希
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
	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "message": err.Error(), "data": nil})
		return
	}

	// 2. 调用 service 层处理登录逻辑
	token, user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 4001, "message": err.Error(), "data": nil})
		return
	}

	// 3. 登录成功，返回 token 和用户信息
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

// GetCurrentUser 获取当前登录用户的信息
// 这个处理器依赖于 JWTMiddleware 将用户信息放入 context
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// 从 context 中获取由中间件设置的 claims
	claims, exists := c.Get("user_claims")
	if !exists {
		// 如果 claims 不存在，说明中间件逻辑有问题，这是一个服务器内部错误
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "无法获取用户信息",
		})
		return
	}

	// 类型断言，将 interface{} 转换为 *service.CustomClaims
	userClaims, ok := claims.(*service.CustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "用户信息格式错误",
		})
		return
	}

	// 这里为了简单，直接返回 claims 中的信息。
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

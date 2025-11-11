// Package middleware 存放 Gin 中间件
package middleware

import (
	"errors"
	"net/http"
	"questflow/internal/service"
	"questflow/pkg/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware 创建一个 Gin 中间件用于 JWT 认证
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    4001,
				"message": "请求未携带 token，无权限访问",
			})
			return
		}

		// 2. 检查 token 格式是否为 "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    4001,
				"message": "token 格式不正确",
			})
			return
		}

		tokenString := parts[1]

		// 3. 解析和验证 token
		token, err := jwt.ParseWithClaims(tokenString, &service.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 确保 token 的签名方法是期望的
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(config.Cfg.JWT.Secret), nil
		})

		// 4. 处理解析错误
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 4001, "message": "token 已过期"})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 4001, "message": "无效的 token"})
			}
			return
		}

		// 5. 验证 claims 并存入 context
		if claims, ok := token.Claims.(*service.CustomClaims); ok && token.Valid {
			// 将 claims 存入 gin.Context，后续的 handler 可以通过 c.Get("user_claims") 来获取
			c.Set("user_claims", claims)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 4001, "message": "无效的 token"})
			return
		}

		// 6. Token 验证通过，继续处理后续请求
		c.Next()
	}
}

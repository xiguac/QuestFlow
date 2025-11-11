// Package handler 存放 HTTP 请求的处理器函数
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 健康检查接口
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "pong",
	})
}

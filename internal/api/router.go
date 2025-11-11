// Package api 负责 API 路由的定义和初始化
package api

import (
	"questflow/internal/api/handler"
	"questflow/internal/api/middleware"
	"questflow/internal/repository"
	"questflow/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 初始化 Gin 引擎并设置所有路由
func SetupRouter(db *gorm.DB) *gin.Engine {
	// 初始化各模块的依赖 (Repository -> Service -> Handler)

	// User
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Submission (需要先于 FormService 初始化，因为 FormService 依赖它)
	submissionRepo := repository.NewSubmissionRepository(db)
	submissionService := service.NewSubmissionService(submissionRepo)

	// Form (依赖 SubmissionRepository)
	formRepo := repository.NewFormRepository(db)
	formService := service.NewFormService(formRepo, submissionRepo)
	formHandler := handler.NewFormHandler(formService)

	// SubmissionHandler (依赖 SubmissionService 和 FormService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, formService)

	// 创建 Gin 引擎
	r := gin.Default()

	// 定义 API 版本 v1 的路由组
	apiV1 := r.Group("/api/v1")
	{
		// --- 公开路由 (不需要认证) ---
		apiV1.GET("/ping", handler.Ping)

		publicRoutes := apiV1.Group("/public")
		{
			publicRoutes.GET("/forms/:form_key", formHandler.GetPublicForm)                       // 获取公开表单
			publicRoutes.POST("/forms/:form_key/submissions", submissionHandler.CreateSubmission) // 提交表单
		}

		userPublicRoutes := apiV1.Group("/users")
		{
			userPublicRoutes.POST("/register", userHandler.Register) // 用户注册
			userPublicRoutes.POST("/login", userHandler.Login)       // 用户登录
		}

		// --- 受保护的路由 (需要 JWT 认证) ---
		authRequired := apiV1.Group("/")
		authRequired.Use(middleware.JWTMiddleware())
		{
			userAuthRoutes := authRequired.Group("/users")
			{
				userAuthRoutes.GET("/me", userHandler.GetCurrentUser) // 获取当前用户信息
			}

			formAuthRoutes := authRequired.Group("/forms")
			{
				formAuthRoutes.POST("/", formHandler.CreateForm)  // 创建新表单
				formAuthRoutes.GET("/my", formHandler.GetMyForms) // 获取当前用户的所有表单

				// 针对特定 form_id 的操作
				formAuthRoutes.GET("/:form_id/stats", formHandler.GetStatistics)     // 获取统计数据
				formAuthRoutes.DELETE("/:form_id", formHandler.DeleteForm)           // 删除表单
				formAuthRoutes.GET("/:form_id/details", formHandler.GetFormDetails)  // 获取详情
				formAuthRoutes.PUT("/:form_id", formHandler.UpdateForm)              // 更新整个表单
				formAuthRoutes.PUT("/:form_id/status", formHandler.UpdateFormStatus) // 更新状态
			}
		}
	}

	return r
}

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
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	submissionRepo := repository.NewSubmissionRepository(db)
	submissionService := service.NewSubmissionService(submissionRepo)
	formRepo := repository.NewFormRepository(db)
	formService := service.NewFormService(formRepo, submissionRepo)
	formHandler := handler.NewFormHandler(formService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, formService)

	r := gin.Default()
	apiV1 := r.Group("/api/v1")
	{
		// --- 公开路由 (不需要认证) ---
		apiV1.GET("/ping", handler.Ping)
		publicRoutes := apiV1.Group("/public")
		{
			publicRoutes.GET("/forms/:form_key", formHandler.GetPublicForm)
			publicRoutes.POST("/forms/:form_key/submissions", submissionHandler.CreateSubmission)
		}
		userPublicRoutes := apiV1.Group("/users")
		{
			userPublicRoutes.POST("/register", userHandler.Register)
			userPublicRoutes.POST("/login", userHandler.Login)
		}

		// --- 受保护的路由 (需要 JWT 认证) ---
		authRequired := apiV1.Group("/")
		authRequired.Use(middleware.JWTMiddleware())
		{
			userAuthRoutes := authRequired.Group("/users")
			{
				userAuthRoutes.GET("/me", userHandler.GetCurrentUser)
				userAuthRoutes.PUT("/profile", userHandler.UpdateProfile)
			}

			formAuthRoutes := authRequired.Group("/forms")
			{
				formAuthRoutes.POST("/", formHandler.CreateForm)
				formAuthRoutes.GET("/my", formHandler.GetMyForms)

				// 针对特定 form_id 的操作
				formAuthRoutes.GET("/:form_id/stats", formHandler.GetStatistics)
				formAuthRoutes.DELETE("/:form_id", formHandler.DeleteForm)
				formAuthRoutes.GET("/:form_id/details", formHandler.GetFormDetails)
				formAuthRoutes.PUT("/:form_id", formHandler.UpdateForm)
				formAuthRoutes.PUT("/:form_id/status", formHandler.UpdateFormStatus)

				// 【核心改动】将导出路由从 GET 修改为 POST
				formAuthRoutes.POST("/:form_id/export", formHandler.ExportSubmissions)
			}
		}
	}
	return r
}

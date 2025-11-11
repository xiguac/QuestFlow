// QuestFlow 项目的唯一启动入口
package main

import (
	"log"
	"questflow/internal/api"
	"questflow/internal/consumer"
	"questflow/internal/model"
	"questflow/pkg/config"
	"questflow/pkg/db"
	"questflow/pkg/redis"
)

func main() {
	// 1. 加载配置
	config.Init("./configs/config.yaml")

	// 2. 初始化数据库连接
	db.InitMySQL()

	// 3. 初始化 Redis 连接
	redis.InitRedis()

	// 4. 自动迁移数据库表结构
	err := db.DB.AutoMigrate(&model.User{}, &model.Form{}, &model.Submission{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	// 5. 在启动 Goroutine
	consumer.StartSubmissionConsumer()

	// 6. 设置并启动 Gin API 服务 (这将阻塞主线程)
	router := api.SetupRouter(db.DB)
	log.Printf("API Server is running on http://localhost%s", config.Cfg.App.Port)
	if err := router.Run(config.Cfg.App.Port); err != nil {
		log.Fatalf("Failed to run API server: %v", err)
	}
}

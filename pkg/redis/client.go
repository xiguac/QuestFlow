// Package redis 负责初始化和管理 Redis 客户端连接
package redis

import (
	"context"
	"log"
	"questflow/pkg/config"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

// InitRedis 初始化 Redis 客户端连接
func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Cfg.Redis.Addr,
		Password: config.Cfg.Redis.Password,
		DB:       config.Cfg.Redis.DB,
	})

	// 测试连接
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection successful.")
}

// PublishSubmissionMessage 封装了向 submission stream 发送消息的逻辑
func PublishSubmissionMessage(ctx context.Context, payload []byte) (string, error) {
	// 使用 redis.XAddArgs 来构造参数
	args := &redis.XAddArgs{
		Stream: config.Cfg.Redis.SubmissionStreamKey,
		Values: map[string]interface{}{"payload": payload},
	}

	// 发送消息到 Redis Stream
	messageID, err := RDB.XAdd(ctx, args).Result()
	if err != nil {
		return "", err
	}

	return messageID, nil
}

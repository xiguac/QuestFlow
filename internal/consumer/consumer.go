// Package consumer 封装了 submission consumer 服务的逻辑
package consumer

import (
	"context"
	"encoding/json"
	"log"
	"questflow/internal/repository"
	"questflow/internal/service"
	"questflow/pkg/config"
	"questflow/pkg/db"
	redisPkg "questflow/pkg/redis"
	"time"

	"github.com/go-redis/redis/v8"
)

// StartSubmissionConsumer 启动 submission consumer
func StartSubmissionConsumer() {
	log.Println("Starting submission consumer goroutine...")

	// 依赖注入
	submissionRepo := repository.NewSubmissionRepository(db.DB)
	submissionService := service.NewSubmissionService(submissionRepo)

	streamKey := config.Cfg.Redis.SubmissionStreamKey
	groupName := config.Cfg.Redis.SubmissionGroupName
	consumerName := "consumer-" + time.Now().Format("20060102150405")

	err := redisPkg.RDB.XGroupCreateMkStream(context.Background(), streamKey, groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	log.Printf("Consumer group '%s' is ready.", groupName)

	go func() {
		for {
			streams, err := redisPkg.RDB.XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{streamKey, ">"},
				Count:    10,
				Block:    0,
			}).Result()

			if err != nil {
				log.Printf("Error reading from stream: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					log.Printf("[Consumer] Processing message ID: %s", message.ID)

					payload, ok := message.Values["payload"].(string)
					if !ok {
						log.Printf("[Consumer] Invalid payload for message %s. Skipping.", message.ID)
						continue
					}

					var subMsg service.SubmissionMessage
					if err := json.Unmarshal([]byte(payload), &subMsg); err != nil {
						log.Printf("[Consumer] Failed to unmarshal message %s: %v. Skipping.", message.ID, err)
						continue
					}

					if err := submissionService.ProcessSubmission(subMsg); err != nil {
						log.Printf("[Consumer] Failed to process message %s: %v", message.ID, err)
						continue
					}

					redisPkg.RDB.XAck(context.Background(), streamKey, groupName, message.ID)
					log.Printf("[Consumer] Successfully processed and ACKed message ID: %s", message.ID)
				}
			}
		}
	}()

	log.Println("Consumer is now listening for messages in the background.")
}

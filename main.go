package main

import (
	"context"
	"control-plane/api"
	"control-plane/controller"
	"control-plane/queue"
	"control-plane/storage"
	"control-plane/worker"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	storageManager := storage.NewRedisManager(REDIS_HOST, REDIS_PORT)
	workerFactory := worker.NewFactory(storageManager)
	redisQueue := queue.NewRedisQueue(queue.RedisQueueConfig{
		RedisHost:          REDIS_HOST,
		RedisPort:          REDIS_PORT,
		RateLimitPerMinute: 1,
	}, storageManager, workerFactory)

	ctrl := controller.NewManager(storageManager, redisQueue)

	v1 := router.Group("/v1")
	v1api := api.NewRouter(ctrl)
	v1api.Router(v1)

	redisQueue.Run(context.Background())
	redisQueue.RegisterWorker()

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("error while starting HTTP server: %v", err)
	}
}

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"ksqldb-trace/api"
	"ksqldb-trace/controller"
	"ksqldb-trace/queue"
	"ksqldb-trace/storage"
	"ksqldb-trace/worker"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	storageManager := storage.NewRedisManager()
	workerFactory := worker.NewFactory(storageManager)
	redisQueue := queue.NewRedisQueue(storageManager, workerFactory)

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

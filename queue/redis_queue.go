package queue

import (
	"context"
	"control-plane/models"
	"control-plane/storage"
	"control-plane/worker"
	"fmt"
	"github.com/adjust/rmq/v5"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"time"
)

type RedisQueueConfig struct {
	RedisHost string
	RedisPort string

	RateLimitPerMinute int
}

type RedisQueue struct {
	RedisQueueConfig

	client     *redis.Client
	connection rmq.Connection
	tasks      rmq.Queue
	manager    storage.Manager
	factory    worker.Factory
}

func NewRedisQueue(config RedisQueueConfig, manager storage.Manager, factory worker.Factory) *RedisQueue {
	redisStr := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	connection, err := rmq.OpenConnection("producer_consumer", "tcp", redisStr, 1, nil)
	if err != nil {
		log.Fatalf("unable to create queue: %v", err)
		return nil
	}

	tasks, err := connection.OpenQueue("generation_task")
	if err != nil {
		connection.StopAllConsuming()
		log.Fatalf("unable to create new queue for managing generation tasks")
		return nil
	}

	return &RedisQueue{RedisQueueConfig: config, connection: connection, tasks: tasks, manager: manager, factory: factory}
}

func (r *RedisQueue) Enqueue(ctx context.Context, generationTask models.GenerationTask) error {
	log.Debugf("Enqueue new generation task %s", generationTask.Id)
	buf, err := generationTask.MarshalBinary()
	if err != nil {
		return err
	}

	r.manager.UpdateGenerationTaskStatus(ctx, generationTask.Id, models.GenerationTaskStatusPending)

	return r.tasks.PublishBytes(buf)
}

func (r *RedisQueue) RegisterWorker() {
	_, err := r.tasks.AddConsumer("worker", r.factory.NewWorker())
	if err != nil {
		log.Fatalf("unable to create consumer for handling generation tasks: %v", err)
	}
}

func (r *RedisQueue) Run(ctx context.Context) {
	log.Info("Starting queue")
	if err := r.tasks.StartConsuming(1000, 100*time.Millisecond); err != nil {
		log.Fatalf("unable to start queue: %v", err)
	}
}

package storage

import (
	"context"
	"control-plane/models"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	GenerationTaskKey       = "generation_tasks"
	GenerationTaskStatusKey = "generation_tasks_status"
)

type RedisManager struct {
	client *redis.Client
}

func (r *RedisManager) StoreGenerationTask(ctx context.Context, task models.GenerationTask, initialStatus models.GenerationTaskStatusType) error {
	if len(task.Id) == 0 {
		return fmt.Errorf("could not store, please provide valid id")
	}

	err := r.client.HSet(ctx, GenerationTaskKey, task.Id, task).Err()
	if err != nil {
		return fmt.Errorf("unable to store generation task %s: %v", task.Id, err)
	}

	return r.UpdateGenerationTaskStatus(ctx, task.Id, initialStatus)
}

func (r *RedisManager) GetGenerationTask(ctx context.Context, id string) (*models.GenerationTask, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("please provide a valid id")
	}

	var generationTask *models.GenerationTask
	if err := r.client.HGet(ctx, GenerationTaskKey, id).Scan(generationTask); err != nil {
		return nil, fmt.Errorf("unable to get generation task for id %s: %v", id, err)
	}
	return generationTask, nil
}

func (r *RedisManager) GetGenerationTaskStatus(ctx context.Context, id string) (*models.GenerationTaskStatus, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("could not get the generation task status, please provide an id")
	}

	var taskStatus models.GenerationTaskStatus
	if err := r.client.HGet(ctx, GenerationTaskStatusKey, id).Scan(&taskStatus); err != nil {
		return nil, fmt.Errorf("unable to load task generation status for id %s: %v", id, err)
	}

	return &taskStatus, nil
}

func (r *RedisManager) UpdateGenerationTaskStatus(ctx context.Context, id string, status models.GenerationTaskStatusType) error {
	if len(id) == 0 {
		return fmt.Errorf("could not update the status, please provide an id")
	}

	var currentTaskStatus *models.GenerationTaskStatus
	currentTaskStatus, _ = r.GetGenerationTaskStatus(ctx, id)
	if currentTaskStatus == nil {
		currentTaskStatus = &models.GenerationTaskStatus{}
		currentTaskStatus.Status = status
		currentTaskStatus.Id = id
		currentTaskStatus.History = make([]models.GenerationTaskStatusHistory, 1)
		currentTaskStatus.History[0].Status = status
		currentTaskStatus.History[0].Timestamp = time.Now()

	} else {
		currentTaskStatus.History = append(currentTaskStatus.History, models.GenerationTaskStatusHistory{
			Status:    status,
			Timestamp: time.Now(),
		})
		currentTaskStatus.Status = status
	}

	return r.client.HSet(ctx, GenerationTaskStatusKey, id, currentTaskStatus).Err()

}

func NewRedisManager(redisHost string, redisPort string) *RedisManager {
	client := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", redisHost, redisPort)})

	if client == nil {
		log.Fatalf("unable to establish connection to redis")
		return nil
	}

	return &RedisManager{
		client: client,
	}
}

func (r *RedisManager) Close() {
	err := r.client.Close()
	if err != nil {
		log.Errorf("error while closing connection: %v", err)
	}
}

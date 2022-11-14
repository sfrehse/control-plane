package controller

import (
	"context"
	"control-plane/models"
	"control-plane/queue"
	"control-plane/storage"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	prefixIdGenerator *PrefixIdGenerator
	storageManager    storage.Manager
	queue             *queue.RedisQueue
}

func NewManager(storageManager storage.Manager, queue *queue.RedisQueue) *Controller {
	idGenerator := NewPrefixedIdGenerator()

	if idGenerator == nil {
		log.Fatalf("unable to initialize controller")
		return nil
	}

	return &Controller{
		prefixIdGenerator: idGenerator,
		storageManager:    storageManager,
		queue:             queue,
	}
}
func (m *Controller) CreateNewTask(ctx context.Context, task models.GenerationTask) (*models.GenerationTaskStatus, error) {
	newId, err := m.prefixIdGenerator.Generator(GenerationTaskIdPrefix)

	if err != nil {
		return nil, fmt.Errorf("unable to create new task since unable to generate new id: %v", err)
	}

	task.Id = newId

	if err := m.storageManager.StoreGenerationTask(ctx, task, models.GenerationTaskStatusOpen); err != nil {
		return nil, fmt.Errorf("unable to store the task: %v", err)
	}

	if err := m.queue.Enqueue(ctx, task); err != nil {
		log.Errorf("unable to enqueue new generation task job %s: %v", task.Id, err)
	}
	log.Debugf("Enqueued new generation task %s", task.Id)

	return m.GetGenerationTaskStatus(ctx, task.Id)
}

func (m *Controller) GetGenerationTaskStatus(ctx context.Context, id string) (*models.GenerationTaskStatus, error) {
	generationTaskStatus, err := m.storageManager.GetGenerationTaskStatus(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get current status for id %s: %v", id, err)
	}

	return generationTaskStatus, nil
}

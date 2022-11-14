package storage

import (
	"context"
	"control-plane/models"
)

type Manager interface {
	StoreGenerationTask(ctx context.Context, task models.GenerationTask, initialStatus models.GenerationTaskStatusType) error
	GetGenerationTask(ctx context.Context, id string) (*models.GenerationTask, error)

	GetGenerationTaskStatus(ctx context.Context, id string) (*models.GenerationTaskStatus, error)
	UpdateGenerationTaskStatus(ctx context.Context, id string, status models.GenerationTaskStatusType) error
}

package storage

import (
	"context"
	"control-plane/models"
)

type Queue interface {
	Enqueue(ctx context.Context, generationTask models.GenerationTask) error
}

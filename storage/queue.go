package storage

import (
	"context"
	"ksqldb-trace/models"
)

type Queue interface {
	Enqueue(ctx context.Context, generationTask models.GenerationTask) error
}

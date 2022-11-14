package storage

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"ksqldb-trace/models"
	"math/rand"
	"testing"
	"time"
)

func TestRedisStorageManager(t *testing.T) {
	ctx := context.Background()

	manager := NewRedisManager()
	assert.NotNil(t, manager)

	generationTask := models.GenerationTask{
		Schema:    "CREATE STREAM INPUT (A INT) WITH (KAFKA_TOPIC='test, VALUE_FORMAT='json');",
		Statement: "CREATE STREAM OUTPUT AS SELECT A * 2 FROM INPUT EMIT CHANGES;",
	}

	err := manager.StoreGenerationTask(ctx, generationTask, models.GenerationTaskStatusOpen)
	assert.Error(t, err)
	rand.Seed(time.Now().Unix())

	generationTask.Id = fmt.Sprintf("random_id_%d", rand.Int63())

	err = manager.StoreGenerationTask(ctx, generationTask, models.GenerationTaskStatusOpen)
	assert.NoError(t, err)

	err = manager.UpdateGenerationTaskStatus(ctx, generationTask.Id, models.GenerationTaskStatusProcessing)

	generationTaskStatus, err := manager.GetGenerationTaskStatus(ctx, generationTask.Id)
	assert.NoError(t, err)

	assert.Equal(t, models.GenerationTaskStatusProcessing, string(generationTaskStatus.Status))
	assert.Equal(t, 2, len(generationTaskStatus.History))
}

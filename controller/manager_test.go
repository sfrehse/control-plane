package controller

import (
	"context"
	"control-plane/models"
	"control-plane/queue"
	"control-plane/storage"
	"control-plane/worker"
	"flag"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var redisHost = flag.String("redis-host", "localhost", "host to redis host")
var redisPort = flag.String("redis-port", "6379", "host to redis port")

func TestManager_CreateNewGenerationTasks(t *testing.T) {
	storageManager := storage.NewRedisManager(*redisHost, *redisPort)
	redisQueue := queue.NewRedisQueue(queue.RedisQueueConfig{
		RedisHost:          *redisHost,
		RedisPort:          *redisPort,
		RateLimitPerMinute: 1,
	}, storageManager, worker.NewFactory(storageManager))

	manager := NewManager(storageManager, redisQueue)
	assert.NotNil(t, manager)

	generationTask := models.GenerationTask{
		Schema:    "CREATE STREAM INPUT (A INT) WITH (KAFKA_TOPIC='test, VALUE_FORMAT='json');",
		Statement: "CREATE STREAM OUTPUT AS SELECT A * 2 FROM INPUT EMIT CHANGES;",
	}

	generationTaskStatus, err := manager.CreateNewTask(context.Background(), generationTask)

	assert.NoError(t, err)
	assert.Equal(t, models.GenerationTaskStatusPending, string(generationTaskStatus.Status))
	assert.True(t, len(generationTaskStatus.Id) > 0)
	assert.True(t, strings.HasPrefix(generationTaskStatus.Id, "qry"))

}

package api

import (
	"bytes"
	"control-plane/controller"
	"control-plane/models"
	"control-plane/queue"
	"control-plane/storage"
	"control-plane/worker"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func SetupRouter() *gin.Engine {
	return gin.Default()
}

func TestRouter_CreateGenerationTask(t *testing.T) {
	storageManager := storage.NewRedisManager()
	workerFactory := worker.NewFactory(storageManager)
	redisQueue := queue.NewRedisQueue(storageManager, workerFactory)

	ctrl := controller.NewManager(storageManager, redisQueue)

	router := SetupRouter()
	NewRouter(ctrl).Router(router.Group("/v1"))

	generationTask := models.GenerationTask{
		Schema:    "CREATE STREAM INPUT (A INT) WITH (KAFKA_TOPIC='test, VALUE_FORMAT='json');",
		Statement: "CREATE STREAM OUTPUT AS SELECT A * 2 FROM INPUT EMIT CHANGES;",
	}

	buf, err := generationTask.MarshalBinary()
	assert.NoError(t, err)
	request := httptest.NewRequest("POST", "/v1/generationTask", bytes.NewBuffer(buf))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	responseData, _ := ioutil.ReadAll(w.Body)

	var generationTaskStatus models.GenerationTaskStatus
	assert.NoError(t, generationTaskStatus.UnmarshalBinary(responseData))

	assert.Equal(t, models.GenerationTaskStatusPending, string(generationTaskStatus.Status))
	assert.True(t, len(generationTaskStatus.Id) > 0)
	assert.True(t, len(generationTaskStatus.History) > 0)
}

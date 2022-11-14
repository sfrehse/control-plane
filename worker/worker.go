package worker

import (
	"context"
	"control-plane/models"
	"control-plane/storage"
	"github.com/adjust/rmq/v5"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	manager storage.Manager
}

func NewWorker(manager storage.Manager) *Worker {
	return &Worker{manager: manager}
}

func (worker *Worker) Consume(delivery rmq.Delivery) {
	var data models.GenerationTask
	if err := data.UnmarshalBinary([]byte(delivery.Payload())); err != nil {
		log.Fatalf("unable to unpack data: %v", err)
		delivery.Reject()
		return
	}

	ctx := context.Background()

	worker.manager.UpdateGenerationTaskStatus(ctx, data.Id, models.GenerationTaskStatusProcessing)
	log.Infof("Processing %v", data.Id)

	worker.manager.UpdateGenerationTaskStatus(ctx, data.Id, models.GenerationTaskStatusCompleted)
	delivery.Ack()
}

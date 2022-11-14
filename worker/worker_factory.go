package worker

import (
	"github.com/adjust/rmq/v5"
	"ksqldb-trace/storage"
)

type Factory interface {
	NewWorker() rmq.Consumer
}

type FactoryImpl struct {
	manager storage.Manager
}

func NewFactory(manager storage.Manager) *FactoryImpl {
	return &FactoryImpl{manager: manager}
}

func (factory *FactoryImpl) NewWorker() rmq.Consumer {
	return NewWorker(factory.manager)
}

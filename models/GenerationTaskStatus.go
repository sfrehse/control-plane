package models

import (
	"encoding/json"
	"time"
)

const (
	GenerationTaskStatusOpen       = "Open"
	GenerationTaskStatusPending    = "Pending"
	GenerationTaskStatusProcessing = "Processing"
	GenerationTaskStatusCompleted  = "Completed"
	GenerationTaskStatusErrored    = "Errored"
)

type GenerationTaskStatusType string

type GenerationTaskStatusHistory struct {
	Timestamp time.Time                `json:"timestamp"`
	Status    GenerationTaskStatusType `json:"status"`
}

type GenerationTaskStatus struct {
	Id     string                   `json:"id"`
	Status GenerationTaskStatusType `json:"status"`

	History []GenerationTaskStatusHistory `json:"history"`
}

func (status GenerationTaskStatus) MarshalBinary() ([]byte, error) {
	return json.Marshal(status)
}

func (status *GenerationTaskStatus) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, status)
}

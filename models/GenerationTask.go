package models

import "encoding/json"

type GenerationTask struct {
	Id        string `json:"id"`
	Schema    string `json:"schema"`
	Statement string `json:"statement"`
}

func (task GenerationTask) MarshalBinary() ([]byte, error) {
	return json.Marshal(task)
}

func (task *GenerationTask) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, task)
}

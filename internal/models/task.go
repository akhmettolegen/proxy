package models

import "time"

const (
	TaskStatusNew       = "new"
	TaskStatusInProcess = "in_process"
	TaskStatusDone      = "done"
	TaskStatusError     = "error"
)

type Task struct {
	Id             string              `bson:"id" json:"id"`
	Status         string              `bson:"status" json:"status"`
	HttpStatusCode int                 `bson:"httpStatusCode" json:"httpStatusCode"`
	Headers        map[string][]string `json:"headers"`
	Length         int                 `bson:"length" json:"length"`
	CreatedAt      time.Time           `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt      time.Time           `bson:"updatedAt" json:"updatedAt,omitempty"`
}

type TaskRequest struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

type TaskResponse struct {
	Id string `json:"id"`
}

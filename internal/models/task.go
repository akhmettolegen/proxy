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
	CreatedAt      time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time           `bson:"updatedAt" json:"updatedAt"`
}

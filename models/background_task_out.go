package models

import "time"

type BackgroundTaskStatus string
const (
	BackgroundTaskStatusRunning  BackgroundTaskStatus = "running"
	BackgroundTaskStatusFinished BackgroundTaskStatus = "finished"
	BackgroundTaskStatusFailed   BackgroundTaskStatus = "failed"
)

type BackgroundTaskType string
const (
	BackgroundTaskTypeApplication BackgroundTaskType = "application"
	BackgroundTaskTypeEndpoint    BackgroundTaskType = "endpoint"
	BackgroundTaskTypeMessage     BackgroundTaskType = "message"
	BackgroundTaskTypeReplay      BackgroundTaskType = "replay"
)

type BackgroundTaskOut struct {
	Id        string               `json:"id"`
	Status    BackgroundTaskStatus `json:"status"`
	Task      BackgroundTaskType   `json:"task"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
}

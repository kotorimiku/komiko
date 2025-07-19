package model

import "time"

type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskPaused    TaskStatus = "paused"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
)

type Task struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Status    TaskStatus `json:"status"`
	Progress  int        `json:"progress"`
	Result    string     `json:"result"`
	Error     string     `json:"error"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

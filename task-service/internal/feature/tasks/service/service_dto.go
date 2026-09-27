package tasks_service

import "time"

type ServiceEventDTO struct {
	TaskID      int       `json:"task_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

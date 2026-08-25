package tasks_transport_http

import (
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
)

type TaskResponseDTO struct {
	ID          int        `json:"id"`
	Version     int        `json:"version"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"comleted_at"`
}

func dtoFromDomain(task *core_domain.Task) TaskResponseDTO {
	return TaskResponseDTO{
		ID:          task.ID,
		Version:     task.Version,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		CompletedAt: task.CompletedAt,
	}
}

func dtoFromDomains(tasks []*core_domain.Task) []TaskResponseDTO {
	response := make([]TaskResponseDTO, len(tasks))
	for i := range tasks {
		dto := dtoFromDomain(tasks[i])
		response[i] = dto
	}

	return response
}

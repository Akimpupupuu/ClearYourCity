package tasks_transport_http

import (
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
)

type TaskResponseDTO struct {
	ID          int
	Version     int
	UserID      int
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	CompletedAt *time.Time
}

func DTOFromDomain(task *core_domain.Task) TaskResponseDTO {
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

func DTOFromDomains(tasks []*core_domain.Task) []TaskResponseDTO {
	response := make([]TaskResponseDTO, len(tasks))
	for i := range tasks {
		dto := DTOFromDomain(tasks[i])
		response[i] = dto
	}

	return response
}

package tasks_transport_http

import (
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
)

type TaskResponseDTO struct {
	ID          int        `json:"id" example:"2"`
	Version     int        `json:"version" example:"1"`
	UserID      int        `json:"user_id" example:"4"`
	Title       string     `json:"title" example:"Мусор"`
	Description string     `json:"description" example:"Мусор около дома по адресу: ул. Пушкина. д. 2"`
	Status      string     `json:"status" example:"created"`
	CreatedAt   time.Time  `json:"created_at" example:"2026-08-29T18:51:08.085831Z"`
	CompletedAt *time.Time `json:"comleted_at" example:"null"`
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

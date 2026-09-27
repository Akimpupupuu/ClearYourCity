package tasks_postgres

import "time"

type TaskEventPayload struct {
	ID          int        `json:"id"`
	Version     int        `json:"version"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Token       string     `json:"token"`
}

func NewTaskEventPayloadFromTaskModel(model TaskModel, token string) TaskEventPayload {
	return TaskEventPayload{
		ID:          model.ID,
		Version:     model.Version,
		UserID:      model.UserID,
		Title:       model.Title,
		Description: model.Description,
		Status:      model.Status,
		CreatedAt:   model.CreatedAt,
		CompletedAt: model.CompletedAt,
		Token:       token,
	}
}

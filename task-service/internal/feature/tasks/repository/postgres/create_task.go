package tasks_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
)

func (r *tasksRepository) CreateTask(ctx context.Context, task *core_domain.Task) (*core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	INSERT 	INTO task_service.task (user_id, title, description, status, created_at, completed_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version, user_id, title, description, status, created_at, completed_at;
	`

	row := r.pool.QueryRow(ctx, query, task.UserID, task.Title, task.Description, string(task.Status), task.CreatedAt, task.CompletedAt)

	var taskModel TaskModel
	if err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.UserID,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Status,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
	); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	taskDomain, err := domainFromModel(taskModel)
	if err != nil {
		return nil, fmt.Errorf("invalid status: %w", err)
	}

	return taskDomain, nil
}

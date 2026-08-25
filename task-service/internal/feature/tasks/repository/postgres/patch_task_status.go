package tasks_postgres

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *tasksRepository) PatchStatus(ctx context.Context, task *core_domain.Task) (*core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	UPDATE tasks_service.tasks
	SET 
		status = $1,
		version = version+1
	WHERE id = $2 AND version = $3
	RETURNING id, version, user_id, title, description, status, created_at, completed_at;
	`

	row := r.pool.QueryRow(ctx, query, task.Status, task.ID, task.Version)

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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task with id = '%d' concurrently accessed: %w", task.ID, core_errors.ErrConflict)
		}

		return nil, fmt.Errorf("scan error: %w", err)
	}

	taskDomain, err := domainFromModel(taskModel)
	if err != nil {
		return nil, fmt.Errorf("create domain from model: %w", err)
	}

	return taskDomain, nil
}

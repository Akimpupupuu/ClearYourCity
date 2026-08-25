package tasks_postgres

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *tasksRepository) GetTask(ctx context.Context, taskID int) (*core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	SELECT id, version, user_id, title, description, status, created_at, completed_at
	FROM tasks_service.task
	WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, taskID)

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
			return nil, fmt.Errorf("task with id = '%d' not found: %w", taskID, core_errors.ErrNotFound)
		}

		return nil, fmt.Errorf("scan error: %w", err)
	}

	taskDomain, err := domainFromModel(taskModel)
	if err != nil {
		return nil, fmt.Errorf("create domain from model: %w", err)
	}

	return taskDomain, nil
}

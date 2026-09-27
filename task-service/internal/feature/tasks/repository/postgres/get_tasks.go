package tasks_postgres

import (
	"context"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
)

func (r *tasksRepository) GetTasks(ctx context.Context, userID int, limit int, offset int) ([]*core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	SELECT id, version, user_id, title, description, status, created_at, completed_at
	FROM task_service.task
	WHERE user_id = $1
	ORDER BY id ASC
	LIMIT $2
	OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	taskModels := make([]TaskModel, 0, limit)

	for rows.Next() {
		var taskModel TaskModel

		err = rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.UserID,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Status,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", rows.Err())
	}

	taskDomains, err := domainsFromModel(taskModels)
	if err != nil {
		return nil, fmt.Errorf("create domain from model: %w", err)
	}

	return taskDomains, nil
}

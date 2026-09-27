package tasks_postgres

import (
	"context"
	"encoding/json"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
)

func (r *tasksRepository) CreateTask(ctx context.Context, task *core_domain.Task, token string) (*core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	taskQuery := `
	INSERT 	INTO task_service.task (user_id, title, description, status, created_at, completed_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version, user_id, title, description, status, created_at, completed_at;
	`

	row := tx.QueryRow(ctx, taskQuery, task.UserID, task.Title, task.Description, string(task.Status), task.CreatedAt, task.CompletedAt)

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
		return nil, fmt.Errorf("create domain from model: %w", err)
	}

	messageQuery := `
	INSERT INTO task_service.message(task_id, payload)
	VALUES ($1, $2);
	`

	dto := NewTaskEventPayloadFromTaskModel(taskModel, token)

	payload, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	_, err = tx.Exec(ctx, messageQuery, taskDomain.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("insert message into repository: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return taskDomain, nil
}

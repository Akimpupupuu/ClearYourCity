package tasks_postgres

import (
	"context"
	"fmt"
)

func (r *tasksRepository) GetPendingMessages(ctx context.Context, limit int) ([]MessageModel, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	query := `
	SELECT id, task_id, status, payload
	FROM task_service.message
	WHERE status = 'pending'
	ORDER BY id ASC
	LIMIT $1
	FOR UPDATE SKIP LOCKED;
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("get rows from repository: %w", err)
	}
	defer rows.Close()

	messages := make([]MessageModel, 0, limit)

	for rows.Next() {
		var message MessageModel
		err = rows.Scan(
			&message.ID,
			&message.TaskID,
			&message.Status,
			&message.Payload,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		messages = append(messages, message)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("next rows: %w", rows.Err())
	}

	return messages, nil
}

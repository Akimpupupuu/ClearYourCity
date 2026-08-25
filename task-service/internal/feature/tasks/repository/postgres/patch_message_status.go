package tasks_postgres

import (
	"context"
	"fmt"
)

func (r *tasksRepository) PatchMessageStatus(ctx context.Context, ids []int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout)
	defer cancel()

	if len(ids) == 0 {
		return nil
	}

	query := `
	UPDATE tasks_service.message
	SET status = 'processed'
	WHERE id = ANY($1);
	`

	_, err := r.pool.Exec(ctx, query, ids)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

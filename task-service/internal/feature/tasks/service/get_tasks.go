package tasks_service

import (
	"context"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
)

func (s *tasksService) GetTasks(ctx context.Context, userID int, limit int, offset int) ([]*core_domain.Task, error) {
	if err := validateLimitOffset(limit, offset); err != nil {
		return nil, fmt.Errorf("validate pagination values: %w", err)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks from repository: %w", err)
	}

	return tasks, nil
}

func validateLimitOffset(limit int, offset int) error {
	const (
		maxOffset = 100
	)
	if limit < 0 {
		return fmt.Errorf("'limit' must be non-negative value: %d: %w", limit, core_errors.ErrInvalidArgument)
	}

	if offset < 0 || offset > maxOffset {
		return fmt.Errorf("'offset' must be non-negative value and smaller then 100: %d: %w", offset, core_errors.ErrInvalidArgument)
	}

	return nil
}

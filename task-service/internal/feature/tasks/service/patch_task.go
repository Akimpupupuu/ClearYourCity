package tasks_service

import (
	"context"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
)

func (s *tasksService) PatchTask(ctx context.Context, taskID int, userID int, patchTaskCommand core_domain.PatchTaskCommand) (*core_domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task from repository: %w", err)
	}

	if task.UserID != userID {
		return nil, fmt.Errorf("access denied: %w", core_errors.ErrNotFound)
	}

	if err = task.ApplyPatch(patchTaskCommand.Title, patchTaskCommand.Description); err != nil {
		return nil, fmt.Errorf("apply patch: %w", err)
	}

	task, err = s.tasksRepository.PatchTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("patch task in repository: %w", err)
	}

	return task, nil
}

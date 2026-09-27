package tasks_service

import (
	"context"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
)

func (s *tasksService) PatchStatus(ctx context.Context, status string, token string) (*core_domain.Task, error) {
	statusDomain, err := core_domain.NewStatus(status)
	if err != nil {
		return nil, fmt.Errorf("validate status: %w", err)
	}

	taskID, err := s.tasksRedis.GetAction(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("get task id from redis: %w", err)
	}

	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task from repository: %w", err)
	}

	if task.Status == statusDomain {
		return task, nil
	}

	task.Status = statusDomain

	task, err = s.tasksRepository.PatchStatus(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("change status in repository: %w", err)
	}

	switch statusDomain {
	case core_domain.StatusInProgress:
		if err := s.tasksRedis.ProlongAction(ctx, token); err != nil {
			s.log.Error("failed to prolong token in redis", core_logger.Err(err))
		}
	case core_domain.StatusDone, core_domain.StatusRejected:
		if err := s.tasksRedis.DeleteRedis(ctx, token); err != nil {
			s.log.Error("failed to delete token in redis", core_logger.Err(err))
		}
	}

	return task, nil
}

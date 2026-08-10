package tasks_service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
)

const (
	notificationTimeout = 5 * time.Second
)

func (s *tasksService) CreateTask(ctx context.Context, task *core_domain.Task) (*core_domain.Task, error) {
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("validate task: %w", err)
	}

	task, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("create task in repository: %w", err)
	}

	event := ServiceEventDTO{
		TaskID:      task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
	}

	detachedCtx := context.WithoutCancel(ctx)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		asyncCtx, cancel := context.WithTimeout(detachedCtx, notificationTimeout)
		defer cancel()

		logger := core_logger.FromContext(asyncCtx)

		token, err := generateCryptoToken()
		if err != nil {
			logger.Warn("generate crypto token", core_logger.Err(err))
			return
		}

		if err = s.tasksRedis.CreateAction(asyncCtx, token, event.TaskID); err != nil {
			logger.Warn("create task action in redis", core_logger.Err(err))
			return
		}

		if err = s.tasksKafka.PublishTaskCreated(asyncCtx, event); err != nil {
			logger.Warn("produce message into kafka", core_logger.Err(err))
			return
		}
	}()

	return task, nil
}

func generateCryptoToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("crypto rand read: %w", err)
	}

	token := hex.EncodeToString(b)
	return token, nil
}

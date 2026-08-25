package tasks_service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
)

func (s *tasksService) CreateTask(ctx context.Context, task *core_domain.Task) (*core_domain.Task, error) {
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("validate task: %w", err)
	}

	token, err := generateCryptoToken()
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	task, err = s.tasksRepository.CreateTask(ctx, task, token)
	if err != nil {
		return nil, fmt.Errorf("create task in repository: %w", err)
	}

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

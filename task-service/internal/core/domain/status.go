package core_domain

import (
	"fmt"

	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
)

type status string

const (
	StatusCreated    status = "created"
	StatusInProgress status = "in_progress"
	StatusDone       status = "done"
	StatusRejected   status = "rejected"
)

func NewStatus(value string) (status, error) {
	class := status(value)
	switch class {
	case StatusCreated, StatusInProgress, StatusDone, StatusRejected:
		return class, nil
	default:
		return "", fmt.Errorf("invalid status: %w", core_errors.ErrInvalidArgument)
	}
}

func NewStatusCreated() status {
	return status(StatusCreated)
}

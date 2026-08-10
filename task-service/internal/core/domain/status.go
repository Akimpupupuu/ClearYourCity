package core_domain

import (
	"fmt"

	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
)

type status string

const (
	statusCreated    status = "created"
	statusInProgress status = "in_progress"
	statusDone       status = "done"
	statusRejected   status = "rejected"
)

func NewStatus(value string) (status, error) {
	class := status(value)
	switch class {
	case statusCreated, statusInProgress, statusDone, statusRejected:
		return class, nil
	default:
		return "", fmt.Errorf("invalid status: %w", core_errors.ErrInvalidArgument)
	}
}

func NewStatusCreated() status {
	return status(statusCreated)
}

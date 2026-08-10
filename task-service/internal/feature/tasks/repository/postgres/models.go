package tasks_postgres

import (
	"fmt"
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
)

type TaskModel struct {
	ID          int
	Version     int
	UserID      int
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	CompletedAt *time.Time
}

func domainFromModel(model TaskModel) (*core_domain.Task, error) {
	status, err := core_domain.NewStatus(model.Status)
	if err != nil {
		return nil, fmt.Errorf("create status: %w", err)
	}

	return core_domain.NewTask(
		model.ID,
		model.Version,
		model.UserID,
		model.Title,
		model.Description,
		status,
		model.CreatedAt,
		model.CompletedAt,
	), nil
}

func domainsFromModel(models []TaskModel) ([]*core_domain.Task, error) {
	domains := make([]*core_domain.Task, len(models))

	for i := range models {
		domain, err := domainFromModel(models[i])
		if err != nil {
			return nil, fmt.Errorf("convert model to domain: %w", err)
		}

		domains[i] = domain
	}

	return domains, nil
}

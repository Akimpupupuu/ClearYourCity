package core_domain

import (
	"fmt"
	"time"

	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
)

var (
	UninitializedID      = -1
	UninitializedVersion = -1
)

type Task struct {
	ID          int
	Version     int
	UserID      int
	Title       string
	Description string
	Status      status
	CreatedAt   time.Time
	CompletedAt *time.Time
}

func NewTask(
	id int,
	version int,
	userID int,
	title string,
	description string,
	status status,
	createdAt time.Time,
	completedAt *time.Time,
) *Task {
	return &Task{
		ID:          id,
		Version:     version,
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      status,
		CreatedAt:   createdAt,
		CompletedAt: completedAt,
	}
}

func NewTaskUninitialized(title string, description string, userID int) *Task {
	return &Task{
		ID:          UninitializedID,
		Version:     UninitializedVersion,
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      NewStatusCreated(),
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
}

func (t *Task) Validate() error {
	titleLength := len([]rune(t.Title))
	if titleLength < 3 || titleLength > 100 {
		return fmt.Errorf("invalid 'title' length: %d: %w", titleLength, core_errors.ErrInvalidArgument)
	}

	descriptionLength := len([]rune(t.Description))
	if descriptionLength < 10 || descriptionLength > 1000 {
		return fmt.Errorf("invalid 'description' length: %d: %w", descriptionLength, core_errors.ErrInvalidArgument)
	}

	if (t.Status == StatusCreated || t.Status == StatusInProgress) && t.CompletedAt != nil {
		return fmt.Errorf("'completed_at' must be nil if 'status' = 'created' or 'in_progress': %w", core_errors.ErrInvalidArgument)
	}

	if (t.Status == StatusDone || t.Status == StatusRejected) && t.CompletedAt == nil {
		return fmt.Errorf("'completed_at' can't be nil if 'status' = 'done' or 'rejected': %w", core_errors.ErrInvalidArgument)
	}

	if t.UserID == -1 {
		return fmt.Errorf("'user_id' can't be nil: %w", core_errors.ErrInvalidArgument)
	}

	if t.CompletedAt != nil {
		if t.CreatedAt.After(*t.CompletedAt) {
			return fmt.Errorf("invalid complition time: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

func (t *Task) ApplyPatch(title *string, description *string) error {
	tmp := *t

	if title != nil {
		tmp.Title = *title
	}

	if description != nil {
		tmp.Description = *description
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate task: %w", err)
	}

	*t = tmp
	return nil
}

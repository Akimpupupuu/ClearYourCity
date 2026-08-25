package tasks_service

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
)

type tasksService struct {
	tasksRepository TasksRepository
	tasksRedis      TasksRedis
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task *core_domain.Task, token string) (*core_domain.Task, error)
	GetTasks(ctx context.Context, userID int, limit int, offset int) ([]*core_domain.Task, error)
	GetTask(ctx context.Context, taskID int) (*core_domain.Task, error)
	PatchTask(ctx context.Context, task *core_domain.Task) (*core_domain.Task, error)
	PatchStatus(ctx context.Context, task *core_domain.Task) (*core_domain.Task, error)
}

type TasksRedis interface {
	GetAction(ctx context.Context, token string) (int, error)
}

func NewTasksService(
	tasksRepository TasksRepository,
	tasksRedis TasksRedis,
) *tasksService {
	return &tasksService{
		tasksRepository: tasksRepository,
		tasksRedis:      tasksRedis,
	}
}

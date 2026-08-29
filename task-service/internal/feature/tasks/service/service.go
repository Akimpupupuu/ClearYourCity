package tasks_service

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
)

type tasksService struct {
	tasksRepository TasksRepository
	tasksRedis      TasksRedis
	log             core_logger.Logger
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
	ProlongAction(ctx context.Context, token string) error
	DeleteRedis(ctx context.Context, token string) error
}

func NewTasksService(
	tasksRepository TasksRepository,
	tasksRedis TasksRedis,
	log core_logger.Logger,
) *tasksService {
	return &tasksService{
		tasksRepository: tasksRepository,
		tasksRedis:      tasksRedis,
		log:             log,
	}
}

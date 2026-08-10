package tasks_service

import (
	"context"
	"sync"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
)

type tasksService struct {
	tasksRepository TasksRepository
	tasksRedis      TasksRedis
	tasksKafka      TasksKafka
	wg              sync.WaitGroup
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task *core_domain.Task) (*core_domain.Task, error)
	GetTasks(ctx context.Context, userID int, limit int, offset int) ([]*core_domain.Task, error)
}

type TasksRedis interface {
	CreateAction(ctx context.Context, token string, taskID int) error
}

type TasksKafka interface {
	PublishTaskCreated(ctx context.Context, event ServiceEventDTO) error
}

func NewTasksService(
	tasksRepository TasksRepository,
	tasksRedis TasksRedis,
	tasksKafka TasksKafka,
) *tasksService {
	return &tasksService{
		tasksRepository: tasksRepository,
		tasksRedis:      tasksRedis,
		tasksKafka:      tasksKafka,
	}
}

func (s *tasksService) GracefulShutdown() {
	s.wg.Wait()
}

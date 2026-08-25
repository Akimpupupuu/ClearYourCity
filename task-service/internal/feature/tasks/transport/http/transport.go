package tasks_transport_http

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	"github.com/go-chi/chi"
)

type tasksHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(ctx context.Context, task *core_domain.Task) (*core_domain.Task, error)
	GetTasks(ctx context.Context, userID int, limit int, offset int) ([]*core_domain.Task, error)
	PatchTask(ctx context.Context, taskID int, userID int, patchTaskCommand core_domain.PatchTaskCommand) (*core_domain.Task, error)
	PatchStatus(ctx context.Context, status string, token string) (*core_domain.Task, error)
}

func NewTasksHandler(tasksService TasksService) *tasksHandler {
	return &tasksHandler{
		tasksService: tasksService,
	}
}

func (h *tasksHandler) Register(router chi.Router) {
	// удаление задачи
	// регистрация маршрутов
}

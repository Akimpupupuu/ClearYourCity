package tasks_transport_http

import (
	"context"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_jwt "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/jwt"
	http_middleware "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/middleware"
	"github.com/go-chi/chi"
)

type tasksHandler struct {
	tasksService   TasksService
	tokenGenerator *core_jwt.TokenGenerator
}

type TasksService interface {
	CreateTask(ctx context.Context, task *core_domain.Task) (*core_domain.Task, error)
	GetTasks(ctx context.Context, userID int, limit int, offset int) ([]*core_domain.Task, error)
	PatchTask(ctx context.Context, taskID int, userID int, patchTaskCommand core_domain.PatchTaskCommand) (*core_domain.Task, error)
	PatchStatus(ctx context.Context, status string, token string) (*core_domain.Task, error)
}

func NewTasksHandler(tasksService TasksService, tokenGenerator *core_jwt.TokenGenerator) *tasksHandler {
	return &tasksHandler{
		tasksService:   tasksService,
		tokenGenerator: tokenGenerator,
	}
}

func (h *tasksHandler) Register(router chi.Router) {
	router.Route("/task", func(subRouter chi.Router) {
		subRouter.Patch("/status", h.PatchStatus)

		subRouter.Group(func(protected_router chi.Router) {
			protected_router.Use(http_middleware.Auth(h.tokenGenerator))

			protected_router.Post("/", h.CreateTask)
			protected_router.Get("/", h.GetTasks)
			protected_router.Patch("/{id}", h.PatchTask)
		})
	})
}

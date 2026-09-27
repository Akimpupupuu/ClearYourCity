package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_domain "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
	core_jwt "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/jwt"
	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
	http_request "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/request"
	http_response "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100" example:"Мусор"`
	Description string `json:"description" validate:"required,min=10,max=1000" example:"Мусор около дома по адресу: ул. Пушкина. д. 2"`
}

type CreateTaskResponse TaskResponseDTO

// CreateTask 	godoc
// @Summary 	Create task
// @Description Create new task in our system
// @Tags 		task
// @Accept 		json
// @Produce 	json
// @Param 		request body CreateTaskRequest true "CreateTask request body"
// @Success 	201 {object} CreateTaskResponse "Succesfully created task"
// @Failure 	400 {object} http_response.ErrorResponse "Bad request"
// @Failure 	401 {object} http_response.ErrorResponse "Unauthorized"
// @Failure 	500 {object} http_response.ErrorResponse "Internal server error"
// @Security    Auth
// @Router 		/task [post]
func (h *tasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	claims, ok := core_jwt.FromContext(ctx)
	if !ok {
		err := fmt.Errorf("get user claims from context: %w", core_errors.ErrUnauthorized)
		responseHandler.ErrorResponse(err, "failed to get task")
		return
	}

	var request CreateTaskRequest
	if err := http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "decode and validate HTTP request")
		return
	}

	domain := domainFromDTO(request, claims.UserID)
	task, err := h.tasksService.CreateTask(ctx, domain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")
		return
	}

	response := CreateTaskResponse(dtoFromDomain(task))
	responseHandler.JsonResponse(response, http.StatusCreated)
}

func domainFromDTO(request CreateTaskRequest, userID int) *core_domain.Task {
	return core_domain.NewTaskUninitialized(request.Title, request.Description, userID)
}

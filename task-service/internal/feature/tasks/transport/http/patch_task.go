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

const taskIDPathValue = "id"

type PatchTaskRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=3,max=100" example:"Мусор"`
	Description *string `json:"description" validate:"omitempty,min=10,max=1000" example:"Мусор около дома по адресу: ул. Пушкина. д. 2"`
}

type PatchTaskResponse TaskResponseDTO

// PatchTask 	godoc
// @Summary 	Patch task
// @Description Patch task in our system
// @Tags 		task
// @Accept 		json
// @Produce 	json
// @Param 		request body PatchTaskRequest true "PatchTask request body"
// @Param      	id path int true "ID of the patching task"
// @Success 	200 {object} PatchTaskResponse "Succesfully patched task"
// @Failure 	400 {object} http_response.ErrorResponse "Bad request"
// @Failure 	401 {object} http_response.ErrorResponse "Unauthorized"
// @Failure 	404 {object} http_response.ErrorResponse "Not found"
// @Failure 	409 {object} http_response.ErrorResponse "Conflict"
// @Failure 	500 {object} http_response.ErrorResponse "Internal server error"
// @Security    Auth
// @Router 		/task/{id} [patch]
func (h *tasksHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	claims, ok := core_jwt.FromContext(ctx)
	if !ok {
		err := fmt.Errorf("get user claims from context: %w", core_errors.ErrUnauthorized)
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	taskID, err := http_request.GetIntPathValue(r, taskIDPathValue)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	var request PatchTaskRequest
	if err := http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "decode and validate HTTP request")
		return
	}

	patchTaskCommand := core_domain.NewPatchTaskCommand(request.Title, request.Description)

	task, err := h.tasksService.PatchTask(ctx, taskID, claims.UserID, patchTaskCommand)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchTaskResponse(dtoFromDomain(task))
	responseHandler.JsonResponse(response, http.StatusOK)
}

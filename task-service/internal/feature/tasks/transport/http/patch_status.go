package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
	http_request "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/request"
	http_response "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/response"
)

func (h *tasksHandler) PatchStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	status, token, err := getStatusTokenQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "query params are empty")
		return
	}

	task, err := h.tasksService.PatchStatus(ctx, status, token)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch status")
		return
	}

	response := dtoFromDomain(task)
	responseHandler.JsonResponse(response, http.StatusOK)
}

func getStatusTokenQueryParams(r *http.Request) (string, string, error) {
	const (
		statusKey = "status"
		tokenKey  = "token"
	)

	status := http_request.GetStringQueryParam(r, statusKey)
	if status == nil {
		return "", "", fmt.Errorf("'status' query param can't be empty: %w", core_errors.ErrInvalidArgument)
	}

	token := http_request.GetStringQueryParam(r, tokenKey)
	if token == nil {
		return "", "", fmt.Errorf("'token' query param can't be empty: %w", core_errors.ErrInvalidArgument)
	}

	return *status, *token, nil
}

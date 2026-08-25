package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
	core_jwt "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/jwt"
	core_logger "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/logger"
	http_request "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/request"
	http_response "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/transport/http/response"
)

func (h *tasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	claims, ok := core_jwt.FromContext(ctx)
	if !ok {
		err := fmt.Errorf("get user claims from context: %w", core_errors.ErrUnauthorized)
		responseHandler.ErrorResponse(err, "failed to get task")
		return
	}

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query params")
		return
	}

	tasks, err := h.tasksService.GetTasks(ctx, claims.UserID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	response := dtoFromDomains(tasks)
	responseHandler.JsonResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (int, int, error) {
	const (
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
		defaultLimit     = 20
		defaultOffset    = 0
	)

	limit, err := http_request.GetIntQueryParam(r, limitQueryParam)
	if err != nil {
		return 0, 0, fmt.Errorf("get limit query param: %w", err)
	}

	var newLimit int
	if limit == nil {
		newLimit = defaultLimit
	} else {
		newLimit = *limit
	}

	offset, err := http_request.GetIntQueryParam(r, offsetQueryParam)
	if err != nil {
		return 0, 0, fmt.Errorf("get offset query param: %w", err)
	}

	var newOffset int
	if offset == nil {
		newOffset = defaultOffset
	} else {
		newOffset = *offset
	}

	return newLimit, newOffset, nil
}

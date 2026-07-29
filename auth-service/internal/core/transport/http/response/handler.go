package http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	"go.uber.org/zap"
)

type ResponseHandler struct {
	log *zap.Logger
	w   http.ResponseWriter
}

func NewResponseHandler(log *zap.Logger, w http.ResponseWriter) *ResponseHandler {
	return &ResponseHandler{
		log: log,
		w:   w,
	}
}

func (h *ResponseHandler) JsonResponse(responseBody any, statusCode int) {
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *ResponseHandler) NoContentResponse() {
	h.w.WriteHeader(http.StatusNoContent)
}

func (h *ResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(msg string, fields ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, core_errors.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))
	h.errorResponse(statusCode, err, msg)
}

func (h *ResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.errorResponse(statusCode, err, msg)
}

func (h *ResponseHandler) errorResponse(statusCode int, err error, msg string) {
	response := ErrorResponse{
		Message: msg,
		Error:   err.Error(),
	}

	h.JsonResponse(response, statusCode)
}

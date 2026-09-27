package http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
	"github.com/go-chi/chi"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := chi.URLParam(r, key)
	if pathValue == "" {
		return 0, fmt.Errorf("no key = '%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf("path value = '%s' by key = '%s' is not a valid integer: %v: %w", pathValue, key, err, core_errors.ErrInvalidArgument)
	}

	return val, nil
}

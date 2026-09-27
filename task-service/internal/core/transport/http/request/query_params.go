package http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/Akimpupupuu/ClearYourCity/task-service/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf("param: '%s', by key: '%s' is not a valid integer: %v: %w", param, key, err, core_errors.ErrInvalidArgument)
	}

	return &val, nil
}

func GetStringQueryParam(r *http.Request, key string) *string {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}

	return &param
}

package users_transport_http

import (
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	http_response "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/response"
)

func (h *usersHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		err = fmt.Errorf("failed to get token from context: %v: %w", err, core_errors.ErrUnauthorized)
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	serviceResponse, err := h.usersService.RefreshToken(ctx, cookie.Value)
	if err != nil {
		if errors.Is(err, core_errors.ErrUnauthorized) {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				Path:     "/api/v1/auth",
			})
		}
		responseHandler.ErrorResponse(err, "failed to refresh token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    serviceResponse.RefreshToken,
		Expires:  serviceResponse.RefreshTokenExpiresAt,
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/v1/auth",
	})

	response := LoginDTOFromService(serviceResponse.AccessToken, serviceResponse.AccessTokenExpiresAt)
	responseHandler.JsonResponse(response, http.StatusOK)

}

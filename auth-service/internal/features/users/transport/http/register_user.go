package users_transport_http

import (
	"errors"
	"net/http"
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	http_request "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/http/request"
	http_response "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/http/response"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	"go.uber.org/zap"
)

type RegisterResponse ResponseRegisterDTO

func (h *UsersHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	// make validation incoming HTTP request
	var request RegisterUserRequest
	if err := http_request.Decode(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode HTTP request")
		return
	}

	userCommand := CommandFromDTO(request)

	serviceResponse, err := h.usersService.RegisterUser(ctx, userCommand)
	if err != nil {
		if errors.Is(err, core_errors.ErrConflict) {
			responseHandler.ErrorResponse(err, "failed to create user")
			return
		}

		if serviceResponse.User.ID != core_domain.UninitializedID {
			log.Warn(
				"user created but session generation failed",
				zap.Error(err),
				zap.Any("user_id", serviceResponse.User.ID),
			)

			response := RegisterResponse(RegisterDTOFromDomain(serviceResponse.User, "", time.Time{}))

			responseHandler.JsonResponse(response, http.StatusAccepted)
			return
		}

		responseHandler.ErrorResponse(err, "failed to create and register user")
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

	response := RegisterResponse(RegisterDTOFromDomain(serviceResponse.User, serviceResponse.AccessToken, serviceResponse.AccessTokenExpiresAt))
	responseHandler.JsonResponse(response, http.StatusCreated)
}

func CommandFromDTO(request RegisterUserRequest) core_domain.RegisterCommand {
	return core_domain.NewCommand(request.FullName, request.Email, request.Password)
}

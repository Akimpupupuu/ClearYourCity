package users_transport_http

import (
	"errors"
	"net/http"
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	http_request "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/request"
	http_response "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/response"
	"go.uber.org/zap"
)

type RegisterUserRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100" example:"Иван Иванов"`
	Email    string `json:"email" validate:"required,min=5,max=100" example:"ivan@gmail.com"`
	Password string `json:"password" validate:"required,min=8" example:"ivan1234"`
}

type RegisterResponse ResponseRegisterDTO

// RegisterUser  godoc
// @Summary 	 Register user
// @Description  Register user in our system
// @Tags 		 user
// @Accept		 json
// @Produce 	 json
// @Param 		 request body RegisterUserRequest true "Register user request body"
// @Success 	 201 {object} RegisterResponse "Succesfully registered user"
// @Failure 	 400 {object} http_response.ErrorResponse "Bad request"
// @Failure 	 404 {object} http_response.ErrorResponse "Not found"
// @Failure 	 409 {object} http_response.ErrorResponse "Conflict"
// @Failure 	 500 {object} http_response.ErrorResponse "Internal server error"
// @Router 		 /auth/register [post]
func (h *usersHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	var request RegisterUserRequest
	if err := http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	registerCommand := registerCommandFromDTO(request)

	serviceResponse, err := h.usersService.RegisterUser(ctx, registerCommand)
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

			response := RegisterResponse(RegisterDTOFromService(serviceResponse.User, "", time.Time{}))
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

	response := RegisterResponse(RegisterDTOFromService(serviceResponse.User, serviceResponse.AccessToken, serviceResponse.AccessTokenExpiresAt))
	responseHandler.JsonResponse(response, http.StatusCreated)
}

func registerCommandFromDTO(request RegisterUserRequest) core_domain.RegisterCommand {
	return core_domain.NewRegisterCommand(request.FullName, request.Email, request.Password)
}

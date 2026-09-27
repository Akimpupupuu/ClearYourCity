package users_transport_http

import (
	"errors"
	"net/http"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	http_request "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/request"
	http_response "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/response"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,min=5,max=100" example:"ivan@gmail.com"`
	Password string `json:"password" validate:"required,min=8" example:"ivan1234"`
}

type LoginResponse ResponseLoginDTO

// LoginUser 	godoc
// @Summary 	Log in user
// @Description Log in user in our system
// @Tags 		user
// @Accept 		json
// @Produce 	json
// @Param 		request body LoginRequest true "Login user request body"
// @Success 	200 {object} LoginResponse "Succesfully loged in user"
// @Failure 	400 {object} http_response.ErrorResponse "Bad request"
// @Failure 	401 {object} http_response.ErrorResponse "Unauthorized"
// @Failure 	404 {object} http_response.ErrorResponse "Not found"
// @Failure 	500 {object} http_response.ErrorResponse "Internal server error"
// @Router 		/auth/login [post]
func (h *usersHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	var request LoginRequest
	if err := http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	loginCommand := loginCommandFromDTO(request)

	serviceResponse, err := h.usersService.LoginUser(ctx, loginCommand)
	if err != nil {
		if errors.Is(err, core_errors.ErrUnauthorized) {
			responseHandler.ErrorResponse(err, "invalid email or password")
			return
		}

		responseHandler.ErrorResponse(err, "failed to login user")
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

	response := LoginResponse(LoginDTOFromService(serviceResponse.AccessToken, serviceResponse.AccessTokenExpiresAt))
	responseHandler.JsonResponse(response, http.StatusOK)
}

func loginCommandFromDTO(request LoginRequest) core_domain.LoginCommand {
	return core_domain.NewLoginCommand(request.Email, request.Password)
}

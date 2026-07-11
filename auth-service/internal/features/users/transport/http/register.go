package users_transport_http

import (
	"net/http"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	http_request "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/http/request"
	http_response "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/http/response"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
)

type RegisterUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse ResponseUserDTO

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

	userDomain := DomainFromDTO(request)

	user, err := h.usersService.RegisterUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create and register user")
		return
	}

	response := RegisterResponse(DTOFromDomain(user))
	responseHandler.JsonResponse(response, http.StatusCreated)
}

func DomainFromDTO(request RegisterUserRequest) core_domain.User {
	return core_domain.NewUserUninitialized(request.FullName, request.Email, request.Password)
}

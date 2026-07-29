package users_transport_http

import (
	"net/http"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	http_request "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/request"
	http_response "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/response"
)

type PatchUserRequest struct {
	FullName *string `json:"full_name" validate:"omitempty,min=3,max=100"`
	Email    *string `json:"email" validate:"omitempty,min=5,max=100"`
}

func (h *usersHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	var request PatchUserRequest
	if err := http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	patchUserCommand := core_domain.NewPatchUserCommand(request.FullName, request.Email)

	user, err := h.usersService.PatchUser(ctx, patchUserCommand)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := UserDTOFromDomain(user)
	responseHandler.JsonResponse(response, http.StatusOK)
}

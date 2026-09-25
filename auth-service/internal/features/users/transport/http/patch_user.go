package users_transport_http

import (
	"net/http"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	http_request "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/request"
	http_response "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/response"
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
)

type PatchUserRequest struct {
	FullName *string `json:"full_name" validate:"omitempty,min=3,max=100" example:"Иван Иванов"`
	Email    *string `json:"email" validate:"omitempty,min=5,max=100" example:"ivan@gmail.com"`
}

type PatchUserResponse ResponseUserDTO

// PatchUser godoc
// @Summary 	 Patch user
// @Description  Patch user
// @Tags 		 user
// @Accept 		 json
// @Produce 	 json
// @Param 		 request body PatchUserRequest true "Patch user request body"
// @Success 	 200 {object} PatchUserResponse "Succesfully patched user"
// @Failure 	 400 {object} http_response.ErrorResponse "Bad request"
// @Failure 	 401 {object} http_response.ErrorResponse "Unauthorized"
// @Failure 	 409 {object} http_response.ErrorResponse "Conflict"
// @Failure 	 500 {object} http_response.ErrorResponse "Internal server error"
// @Security     Auth
// @Router 		 /auth/patch_user [patch]
func (h *usersHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	claims, ok := sessions_jwt.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "failed to get token claims")
		return
	}

	var request PatchUserRequest
	if err := http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	patchUserCommand := core_domain.NewPatchUserCommand(request.FullName, request.Email)

	user, err := h.usersService.PatchUser(ctx, claims.UserID, patchUserCommand)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(UserDTOFromDomain(user))
	responseHandler.JsonResponse(response, http.StatusOK)
}

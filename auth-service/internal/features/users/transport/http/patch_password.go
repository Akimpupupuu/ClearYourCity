package users_transport_http

import (
	"net/http"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	http_request "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/request"
	http_response "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/response"
)

type PatchPasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

func (h *usersHandler) PatchPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	var request PatchPasswordRequest
	if err := http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	patchPasswordCommand := core_domain.NewPatchPasswordCommand(request.OldPassword, request.NewPassword)

	if err := h.usersService.PatchPassword(ctx, patchPasswordCommand); err != nil {
		responseHandler.ErrorResponse(err, "failed to patch password")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/v1/auth",
	})

	responseHandler.NoContentResponse()
}

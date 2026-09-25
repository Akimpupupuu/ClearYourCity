package users_transport_http

import (
	"net/http"

	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	http_response "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/response"
)

// LogoutUser 	godoc
// @Summary 	Log out user
// @Description Log out user from the system
// @Tags 		user
// @Success 	204 "Succesfully loged out user"
// @Failure 	401 {object} http_response.ErrorResponse "Unauthorized"
// @Failure 	500 {object} http_response.ErrorResponse "Internal server error"
// @Router 		/auth/logout [post]
func (h *usersHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "failed to log out user")
		return
	}

	if err := h.usersService.LogoutUser(ctx, cookie.Value); err != nil {
		responseHandler.ErrorResponse(err, "failed to logout user")
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

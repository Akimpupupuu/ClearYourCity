package users_transport_http

import (
	"net/http"

	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
	core_logger "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/logger"
	http_response "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/transport/http/response"
	sessions_jwt "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/features/sessions/jwt"
)

func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := http_response.NewResponseHandler(log, w)

	userClaims, ok := sessions_jwt.FromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "failed to get token claims")
		return
	}

	userDomain, err := h.usersService.GetUser(ctx, userClaims.UserID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	response := UserDTOFromDomain(userDomain)
	responseHandler.JsonResponse(response, http.StatusOK)
}

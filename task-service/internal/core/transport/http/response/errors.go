package http_response

type ErrorResponse struct {
	Message string `json:"message" example:"full error text"`
	Error   string `json:"error" example:"short human-readable message"`
}

package http_response

type ErrorResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}
